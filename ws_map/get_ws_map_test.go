package wsmap

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/redhat-appstudio/segment-bridge.git/kwok"
	"github.com/redhat-appstudio/segment-bridge.git/scripts"
	"github.com/stretchr/testify/assert"
)

func TestGetWSMap(t *testing.T) {
	containerfixture.WithServiceContainer(t, kwok.KwokServiceManifest, func(deployment containerfixture.FixtureInfo) {
		kwok.SetKubeconfig()
		scriptPath := "../scripts/get-workspace-map.sh"

		testCases := []struct {
			name     string
			contexts string
			expected map[string]interface{}
		}{
			{
				"All contexts, all tenant namespaces",
				"kwok-rh01 kwok-m01",
				map[string]interface{}{
					"date-masamune-tenant": "date-masamune",
					"koyasu-tenant":        "koyasu",
					"nobu-tenant":          "nobu",
					"ieyasu-tenant":        "ieyasu",
				},
			},
			{
				"Only rh-member context, no non-tenant namespaces",
				"kwok-rh01",
				map[string]interface{}{
					"ieyasu-tenant": "ieyasu",
					"nobu-tenant":   "nobu",
				},
			},
			{
				"Only member context, no non-tenant namespaces",
				"kwok-m01",
				map[string]interface{}{
					"date-masamune-tenant": "date-masamune",
					"koyasu-tenant":        "koyasu",
					"nobu-tenant":          "nobu",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				os.Setenv("CONTEXTS", tc.contexts)
				output := scripts.AssertExecuteScript(t, scriptPath)
				var outputMap map[string]interface{}
				err := json.Unmarshal(output, &outputMap)
				assert.NoError(t, err, "failed to parse json")
				assert.Equal(t, tc.expected, outputMap, "Test case: %s", tc.name)
			})
		}
	})
}
