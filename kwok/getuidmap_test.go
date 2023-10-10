package kwok

import (
	_ "embed"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/stretchr/testify/assert"
)

//go:embed kwok_container_template.yml
var kwokServiceManifest string

func TestGetUIDMap(t *testing.T) {
	containerfixture.WithServiceContainer(t, kwokServiceManifest, func(deployment containerfixture.FixtureInfo) {
		m, err := setUpClusterConfiguration()
		assert.NoError(t, err, "Failed to set up cluster configuration")

		testCases := []struct {
			name   string
			input  map[string]int64
			expect bool
		}{
			{"Usersignup validation", m, true}, // using 'm' returned from setUpClusterConfiguration
			{"Correct validation", map[string]int64{"petta": 12345678, "blee1": 87654321}, true},
			{"empty key", map[string]int64{"": 12345678}, false},
			{"no value key", map[string]int64{"<no value>": 87654321}, false},
			{"empty map", map[string]int64{}, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expect, validateMap(tc.input), "Test case: %s", tc.name)
			})
		}
	})
}
