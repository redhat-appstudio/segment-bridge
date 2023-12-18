package splunk

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var countPattern = regexp.MustCompile(`(?m)^count\s+(\d+)$`)

func getRecordsCount(splunkAppApiURL, index string) (string, error) {
	SplunkAppSearchURL := GetSearchAPIEndpoint(splunkAppApiURL)
	searchQuery := "search=" + fmt.Sprintf("search index=%s | stats count", index)
	netrcPath, _ := containerfixture.GetNetrcPath()
	cmd := exec.Command("curl", "--netrc-file", netrcPath, SplunkAppSearchURL, "-d", searchQuery, "-d", "output_mode=csv")

	output, err := cmd.Output()

	if err != nil {
		log.Printf("Failed to count %s index records: %v", index, err.Error())
		return "", err
	}

	match := countPattern.FindStringSubmatch(string(output))
	if len(match) < 2 {
		return "", fmt.Errorf("Could not retrieve number of records for index %s", index)
	}
	return string(match[1]), nil
}

func TestGetRecordsCount(t *testing.T) {
	WithSplunkContainer(t, func(deployment containerfixture.FixtureInfo) {
		splunkAppApiURL := GetSplunkAppAPIEndpoint("localhost", deployment.ApiPort, "nobody", "-")
		tests := []struct {
			name       string
			index      string
			want       string
			want_error bool
		}{
			{
				name:  "Simple query with an index that consist of 10 events",
				index: "test_index",
				want:  "10",
			},
			{
				name:  "Simple query with non-existent index should return 0",
				index: "non-existent-index",
				want:  "0",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				records, err := getRecordsCount(splunkAppApiURL, tt.index)
				require.NoError(t, err)
				assert.Equal(t, tt.want, records)
			})
		}
	})
}
