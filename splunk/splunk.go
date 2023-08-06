package splunk

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
)

func getRecordsCount(host, port, owner, app, index string) (string, error) {
	login := fmt.Sprintf("%s:%s", containerfixture.Username, containerfixture.Password)
	endpoint := fmt.Sprintf(
		"https://%s:%s/servicesNS/%s/%s/search/v2/jobs/export?output_mode=csv",
		host, port, owner, app,
	)
	searchQuery := "search=" + fmt.Sprintf("search index=%s | stats count", index)
	cmd := exec.Command("curl", "-u", login, "--insecure", endpoint, "-d", searchQuery)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to count %s index records: %v", index, err)
		return "", err
	}
	result := strings.Split(string(output), "\n")
	return strings.TrimSpace(result[len(result)-2]), nil
}
