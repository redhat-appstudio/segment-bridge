package uidmap

import (
	"encoding/json"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/redhat-appstudio/segment-bridge.git/kwok"
	"github.com/redhat-appstudio/segment-bridge.git/scripts"

	"github.com/stretchr/testify/assert"
)

func TestGetUIDMap(t *testing.T) {
	containerfixture.WithServiceContainer(t, kwok.KwokServiceManifest, func(deployment containerfixture.FixtureInfo) {
		kwok.SetKubeconfig()
		scriptPath := "../scripts/get-uid-map.sh"
		output := scripts.AssertExecuteScript(t, scriptPath)

		var outputMap map[string]interface{}
		err := json.Unmarshal(output, &outputMap)
		assert.NoError(t, err, "failed to parse json")

		assertions := []struct {
			name   string
			input  map[string]interface{}
			expect bool
		}{
			{"Usersignup validation", outputMap, true}, // using 'm' returned from setUpClusterConfiguration
			{"Correct validation", map[string]interface{}{"petta": 12345678, "blee1": 87654321}, true},
			{"empty key", map[string]interface{}{"": 12345678}, false},
			{"no value key", map[string]interface{}{"<no value>": 87654321}, false},
			{"empty map", map[string]interface{}{}, false},
		}

		for _, assertion := range assertions {
			t.Run(assertion.name, func(t *testing.T) {
				assert.Equal(t, assertion.expect, validateUIDMap(assertion.input), "Test case: %s", assertion.name)
			})
		}
	})
}
