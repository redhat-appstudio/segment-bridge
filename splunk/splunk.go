package splunk

import (
	"fmt"
)

const ServiceName = "Splunk"
const ServiceStatusCheckPath = "services/server/status"

// GetSplunkAppAPIEndpoint builds a API URL for a specific Splunk app
func GetSplunkAppAPIEndpoint(host, port, ownerName, appName string) string {
	return fmt.Sprintf("http://%s:%s/servicesNS/%s/%s", host, port, ownerName, appName)
}

// GetSearchAPIEndpoint builds a valid endpoint to be used against a
// Splunk Search API service
func GetSearchAPIEndpoint(appAPIEndpoint string) string {
	return fmt.Sprintf("%s/search/v2/jobs/export", appAPIEndpoint)
}
