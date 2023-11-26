package containerfixture

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/redhat-appstudio/segment-bridge.git/scripts"
)

const (
	serviceUpTimeout       = 120
	containerBuildAttempts = 3
	minDynamicPort         = 49152
	maxDynamicPort         = 65535
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type FixtureInfo struct {
	Yaml    string
	PodName string
	WebPort string
	ApiPort string
}

// Generates a random dynamic port number
func getRandomDynamicPort() string {
	randomPort := fmt.Sprintf("%d", r.Intn(maxDynamicPort-minDynamicPort)+minDynamicPort)
	return randomPort
}

// Modify the yaml manifest to allow running it in parallel
// as many times as needed with dynamic values for its name and ports
func generateFixturePodYaml(templateYaml string) FixtureInfo {
	podName := fmt.Sprintf("test-pod-%s", getRandomDynamicPort())
	webPort := getRandomDynamicPort()
	apiPort := getRandomDynamicPort()
	deployment := FixtureInfo{
		Yaml:    templateYaml,
		PodName: podName,
		WebPort: webPort,
		ApiPort: apiPort,
	}

	template, _ := template.New("manifest").Parse(templateYaml)
	var buf bytes.Buffer
	template.ExecuteTemplate(&buf, "manifest", deployment)
	templatedYAML := buf.String()
	deployment.Yaml = templatedYAML

	return deployment
}

// Building a container image using the 'podman build' command.
func buildContainer(t *testing.T, serviceName, service_image_dir string) {
	cmd := exec.Command("podman", "build", "-t", serviceName, "-f", serviceName+"/Dockerfile")
	cmd.Dir = service_image_dir
	stdout, err := cmd.Output()
	if err == nil {
		t.Logf("STDOUT: %s", string(stdout))
		return
	}
	t.Errorf("Failed to build image")
	t.Errorf("STDERR: %s", err.Error())
	t.Fail()
}

// Building and Running a k8s manifest using 'podman kube play' command.
// if the image is already built, it'll run the pod using it.
// It will make up to three attempts to run with different configurations,
// allowing it to run concurrently if necessary.
func buildAndRunPod(t *testing.T, manifestTemplate string) FixtureInfo {
	var deployment FixtureInfo
	for attempt := 0; attempt < containerBuildAttempts; attempt++ {
		deployment = generateFixturePodYaml(manifestTemplate)
		cmd := exec.Command("podman", "play", "kube", "-")
		cmd.Stdin = strings.NewReader(deployment.Yaml)

		stdout, err := cmd.Output()
		if err == nil {
			t.Logf("STDOUT: %s", string(stdout))
			return deployment
		}
		t.Logf("Failed to start pod (attempt %d of %d).\n", attempt+1, containerBuildAttempts)
		t.Logf("STDERR: %s", err.Error())
	}
	cleanup(t, deployment.PodName)
	t.Fail()
	t.Logf("Failed to start pod after %d attempts", containerBuildAttempts)
	return FixtureInfo{} // Avoid 'missing return' error
}

// RequireServiceIsUp continuously monitors the service container to
// ensure it is up and ready for use in tests. It achieves this by repeatedly
// checking the status endpoint of the service API over a period of two minutes.
func RequireServiceIsUp(t *testing.T, endpoint string, msgAndArgs ...interface{}) {
	netrcPath, _ := GetNetrcPath()
	timeoutStart := time.Now().Unix()
	for {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}",
			"--netrc-file", netrcPath, "--insecure", endpoint)
		output, err := cmd.Output()
		if err != nil && string(output) != "000" {
			t.Log(string(output))
			t.Fatalf("curl command failed: %s", err)
		}
		if string(output) == "200" {
			break
		}
		if time.Now().Unix()-timeoutStart > serviceUpTimeout {
			t.Fatalf("%s", msgAndArgs)
		}
		time.Sleep(5 * time.Second)
	}
}

// GetNetrcPath allows us to get the full path to the containerfixture
// .netrc file from every script in the project
func GetNetrcPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to find the path of the current script")
	}
	dirPath := filepath.Dir(filename)
	netrcPath := filepath.Join(dirPath, ".netrc")

	return netrcPath, nil
}

// Making sure to stop and remove the pod we deployed
func cleanup(t *testing.T, podName string) {
	cmd := exec.Command("podman", "pod", "stop", podName)
	if err := cmd.Run(); err != nil {
		t.Errorf("Error stopping the pod: %s", err)
	}
	removePodCmd := exec.Command("podman", "pod", "rm", "-f", podName)
	if err := removePodCmd.Run(); err != nil {
		t.Fatalf("Error removing the pod: %s", err)
	}
}

func WithServiceContainer(t *testing.T, ServiceManifest string, testFunc func(FixtureInfo)) {
	rootDir, err := scripts.GetRepoRootDir()
	if err != nil {
		t.Errorf("Could not determine path for building the container")
		t.Fail()
	}

	images := getManifestImages(t, ServiceManifest)
	for _, image := range images {
		buildContainer(t, image, rootDir)
	}
	deployment := buildAndRunPod(t, ServiceManifest)

	defer cleanup(t, deployment.PodName)

	testFunc(deployment)
	return
}

type PodDefinition struct {
	Spec struct {
		Containers []struct {
			Image string
		}
	}
}

// Parsing a yaml manifest in order to list the images used in each pod's containers
func getManifestImages(t *testing.T, ServiceManifest string) []string {
	var pod PodDefinition
	if err := yaml.Unmarshal([]byte(ServiceManifest), &pod); err != nil {
		t.Errorf("Error parsing manifest: %v", err)
		t.Fail()
	}

	var images []string
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	return images
}
