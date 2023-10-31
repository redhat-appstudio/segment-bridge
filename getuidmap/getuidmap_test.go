package getuidmap

import (
	_ "embed"
	"os/exec"
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/redhat-appstudio/segment-bridge.git/kwokfixture"
	"github.com/redhat-appstudio/segment-bridge.git/testfixture"
	"github.com/stretchr/testify/assert"
)

//go:embed kwok_container_template.tmpl
var kwokServiceManifest string

const scriptPath = "../scripts/get-uid-map.sh"

type ShellScriptExecutor struct{}

func (s *ShellScriptExecutor) Execute(scriptPath string) ([]byte, error) {
	return exec.Command("/bin/sh", scriptPath).CombinedOutput()
}

func validateMap(m map[string]int64) bool {
	if len(m) == 0 {
		return false
	}

	for user := range m {
		if user == "<no value>" || user == "" {
			return false
		}
	}
	return true
}

func TestGetUIDMap(t *testing.T) {
	containerfixture.WithServiceContainer(t, kwokServiceManifest, func(deployment containerfixture.FixtureInfo) {
		err := kwokfixture.SetUpClusterConfiguration()
		assert.NoError(t, err, "Failed to set up cluster configuration")

		executor := &ShellScriptExecutor{}
		m, err := testfixture.ExecuteAndParseScript(executor, scriptPath)
		assert.NoError(t, err, "ExecuteAndParseScript() should not return an error")

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
