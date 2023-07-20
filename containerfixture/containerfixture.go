package containerfixture

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

//go:embed containerfixture.yaml
var manifest string

const Username = "admin"
const Password = "Password"

// Generates a random dynamic port number
func GetRandomPort() string {
	minPort := 49152
	maxPort := 65535
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomPort := fmt.Sprintf("%d", r.Intn(maxPort-minPort)+minPort)
	return randomPort
}

// Modify the manifest to allow running it in parallel as
// many times as needed with dynamic values
func modifyYaml(yaml string) map[string]string {
	podName := fmt.Sprintf("test-pod-%s", GetRandomPort())
	webPort := GetRandomPort()
	apiPort := GetRandomPort()

	yaml = strings.Replace(yaml, "name: \"test-pod\"", fmt.Sprintf("name: %s", podName), 1)
	yaml = strings.Replace(yaml, "hostPort: 8000", fmt.Sprintf("hostPort: %s", webPort), 1)
	yaml = strings.Replace(yaml, "hostPort: 8089", fmt.Sprintf("hostPort: %s", apiPort), 1)

	return map[string]string{
		"yaml":    yaml,
		"podName": podName,
		"webPort": webPort,
		"apiPort": apiPort,
	}
}

// Building and Running a k8s manifest using 'podman kube play' command.
// It will make up to three attempts to run with different configurations,
// allowing it to run concurrently if necessary
func BuildAndRunPod() map[string]string {
	var deployment map[string]string
	attempts := 3
	for attempt := 0; attempt < attempts; attempt++ {
		deployment := modifyYaml(manifest)

		cmd := exec.Command("podman", "play", "kube", "-")
		cmd.Stdin = strings.NewReader(deployment["yaml"])
		output, err := cmd.CombinedOutput()
		if err == nil {
			outputString := strings.Split(strings.TrimSpace(string(output)), "\n")
			containerID := outputString[len(outputString)-1]
			deployment["containerID"] = containerID
			return deployment
		}
		log.Printf("Failed to start pod (attempt %d of %d).\n", attempt+1, attempts)
	}
	Cleanup("", deployment["podName"])
	log.Fatalf("Failed to start pod after %d attempts", attempts)
	return make(map[string]string) // Avoid 'missing return' error
}

// VerifySplunkServiceIsUp continuously monitors the Splunk container to
// ensure it is up and ready for use in tests. It achieves this by repeatedly
// checking the status endpoint of the Splunk API over a period of two minutes.
func VerifySplunkServiceIsUp(port string) {
	timeoutStart := time.Now().Unix()
	splunkURL := fmt.Sprintf("https://localhost:%s/services/server/status", port)
	for {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "-u",
			fmt.Sprintf("%s:%s", Username, Password), "--insecure", splunkURL)
		output, _ := cmd.CombinedOutput()
		if string(output) == "200" {
			log.Println("Splunk service deployed successfully!")
			break
		}
		if time.Now().Unix()-timeoutStart > 120 {
			log.Fatalf("The Splunk instance is not up, cannot verify indexing for tests.")
		}
		time.Sleep(5 * time.Second)
	}
}

// Making sure to stop and remove the container we deployed
func Cleanup(containerID string, podName string) {
	if containerID != "" {
		log.Println("Stopping and removing the container...")
		cmd := exec.Command("podman", "stop", containerID)
		if err := cmd.Run(); err != nil {
			log.Println("Error stopping the container:", err)
		}
		cmd = exec.Command("podman", "rm", "-f", containerID)
		if err := cmd.Run(); err != nil {
			log.Println("Error removing the container:", err)
		}
	}
	log.Println("Stopping and removing the pod...")
	cmd := exec.Command("podman", "pod", "stop", podName)
	if err := cmd.Run(); err != nil {
		log.Println("Error stopping the pod:", err)
	}
	removePodCmd := exec.Command("podman", "pod", "rm", "-f", podName)
	if err := removePodCmd.Run(); err != nil {
		log.Println("Error removing the pod:", err)
	} else {
		log.Println(fmt.Sprintf("Container %s and Pod %s removed successfully.", containerID, podName))
	}
}
