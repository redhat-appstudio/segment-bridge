package containerfixture

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
	"text/template"
	"time"
)

type DeploymentInfo struct {
	Yaml    string
	PodName string
	WebPort string
	ApiPort string
}

// Generates a random dynamic port number
func getRandomPort() string {
	minPort := 49152
	maxPort := 65535
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomPort := fmt.Sprintf("%d", r.Intn(maxPort-minPort)+minPort)
	return randomPort
}

// Modify the yaml manifest to allow running it in parallel
// as many times as needed with dynamic values
func modifyTemplate(yaml string) DeploymentInfo {

	podName := fmt.Sprintf("test-pod-%s", getRandomPort())
	webPort := getRandomPort()
	apiPort := getRandomPort()
	deployment := DeploymentInfo{
		Yaml:    yaml,
		PodName: podName,
		WebPort: webPort,
		ApiPort: apiPort,
	}

	template, _ := template.New("manifest").Parse(yaml)
	var buf bytes.Buffer
	template.ExecuteTemplate(&buf, "manifest", deployment)
	templatedYAML := buf.String()
	deployment.Yaml = templatedYAML

	return deployment
}

// Building and Running a k8s manifest using 'podman kube play' command.
// It will make up to three attempts to run with different configurations,
// allowing it to run concurrently if necessary
func BuildAndRunPod(t *testing.T, manifest string) DeploymentInfo {
	var deployment map[string]string
	attempts := 3
	for attempt := 0; attempt < attempts; attempt++ {
		deployment := modifyTemplate(manifest)
		cmd := exec.Command("podman", "play", "kube", "-")
		cmd.Stdin = strings.NewReader(deployment.Yaml)
		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err == nil {
			t.Logf("STDOUT: %s", stdoutBuf.String())
			return deployment
		}
		t.Logf("Failed to start pod (attempt %d of %d).\n", attempt+1, attempts)
		t.Logf("STDERR: %s", stderrBuf.String())
	}
	Cleanup(t, deployment["podName"])
	t.Fail()
	t.Logf("Failed to start pod after %d attempts", attempts)
	return DeploymentInfo{"", "", "", ""} // Avoid 'missing return' error
}

// VerifyServiceIsUp continuously monitors the service container to
// ensure it is up and ready for use in tests. It achieves this by repeatedly
// checking the status endpoint of the service API over a period of two minutes.
func VerifyServiceIsUp(t *testing.T, endpoint string, serviceName string) {
	timeoutStart := time.Now().Unix()
	for {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}",
			"--netrc-file", "../containerfixture/.netrc", "--insecure", endpoint)
		output, _ := cmd.CombinedOutput()
		if string(output) == "200" {
			t.Logf("%s service deployed successfully!", serviceName)
			break
		}
		if time.Now().Unix()-timeoutStart > 120 {
			t.Fatalf("The %s instance is not up, cannot verify indexing for tests.", serviceName)
		}
		time.Sleep(5 * time.Second)
	}
}

// Making sure to stop and remove the pod we deployed
func Cleanup(t *testing.T, podName string) {
	t.Log("Stopping and removing the pod...")
	cmd := exec.Command("podman", "pod", "stop", podName)
	if err := cmd.Run(); err != nil {
		t.Log("Error stopping the pod:", err)
	}
	removePodCmd := exec.Command("podman", "pod", "rm", "-f", podName)
	if err := removePodCmd.Run(); err != nil {
		t.Log("Error removing the pod:", err)
	} else {
		t.Logf("Pod %s removed successfully.", podName)
	}
}
