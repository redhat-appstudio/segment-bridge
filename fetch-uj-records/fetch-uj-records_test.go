package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/redhat-appstudio/segment-bridge.git/splunk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	scriptPath   = "../scripts/fetch-uj-records.sh"
	filePathPass = "./requiredOutput"
)

func compareOutputs(t *testing.T, output []byte, filePath string) bool {
	fileOutput, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	return bytes.Equal(fileOutput, output)
}

func runAndValidateScript(t *testing.T, filePath, scriptPath string) bool {
	cmd := exec.Command(scriptPath)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	return compareOutputs(t, output, filePath)
}

func TestFetchUjRecords(t *testing.T) {
	require.NoError(t, os.Setenv("SPLUNK_APP_NAME", "-"), "Failed to set SPLUNK_API_URL")
	require.NoError(t, os.Setenv("SPLUNK_INDEX", "test_index"), "Failed to set SPLUNK_INDEX")
	require.NoError(t, os.Setenv("QUERY_EARLIEST_TIME", "2023-01-01T10:46:34.955712833Z"), "Failed to set QUERY_EARLIEST_TIME")
	splunk.WithSplunkContainer(t, func(deployment containerfixture.FixtureInfo) {
		splunkAppApiURL := splunk.GetSplunkAppAPIEndpoint("localhost", deployment.ApiPort, "nobody", "-")
		require.NoError(t, os.Setenv("SPLUNK_API_URL", "localhost:"+deployment.ApiPort), "Failed to set SPLUNK_API_URL")
		netrcPath, err := containerfixture.GetNetrcPath()
		if err != nil {
			t.Fatalf("error: %s", err.Error())
		}
		require.NoError(t, os.Setenv("CURL_NETRC", netrcPath), "Failed to set CURL_NETRC")
		require.NoError(t, os.Setenv("SPLUNK_APP_API_URL", splunkAppApiURL), "Failed to set SPLUNK_API_URL")
		t.Run("PassPath", func(t *testing.T) {
			assert.True(t, runAndValidateScript(t, filePathPass, scriptPath), "Script validation failed for PassPath")
		})
	})
}
