package splunk

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
)

const ServiceName = "Splunk"
const ServiceStatusCheckPath = "services/server/status"
const NotUpErrorMsg = "The %s instance is not up, cannot verify indexing for tests."

// GetSplunkAppAPIEndpoint builds a API URL for a specific Splunk app
func GetSplunkAppAPIEndpoint(host, port, ownerName, appName string) string {
	return fmt.Sprintf("http://%s:%s/servicesNS/%s/%s", host, port, ownerName, appName)
}

// GetSearchAPIEndpoint builds a valid endpoint to be used against a
// Splunk Search API service
func GetSearchAPIEndpoint(appAPIEndpoint string) string {
	return fmt.Sprintf("%s/search/v2/jobs/export", appAPIEndpoint)
}

//go:embed splunk_container_template.tmpl
var splunkServiceManifest string

func WithSplunkContainer(t *testing.T, testFunc func(containerfixture.FixtureInfo)) {
	containerfixture.WithServiceContainer(
		t, splunkServiceManifest,
		func(fi containerfixture.FixtureInfo) {
			endpoint := fmt.Sprintf("http://localhost:%s/%s", fi.ApiPort, ServiceStatusCheckPath)
			containerfixture.RequireServiceIsUp(t, endpoint, NotUpErrorMsg, ServiceName)
			testFunc(fi)
		})
}
