package splunk

import (
	"testing"

	"github.com/redhat-appstudio/segment-bridge.git/containerfixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRecordsCount(t *testing.T) {
	deployment := containerfixture.BuildAndRunPod()
	containerfixture.VerifySplunkServiceIsUp(deployment["apiPort"])
	defer containerfixture.Cleanup(deployment["containerID"], deployment["podName"])

	tests := []struct {
		name       string
		host       string
		port       string
		index      string
		want       string
		want_error bool
	}{
		{
			name:  "Simple query",
			host:  "localhost",
			port:  deployment["apiPort"],
			index: "test_index",
			want:  "10",
		},
		{
			name:       "Simple query",
			host:       "fake-host",
			port:       deployment["apiPort"],
			index:      "test_index",
			want_error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := getRecordsCount(tt.host, tt.port, "nobody", "-", tt.index)
			if tt.want_error {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, records)
		})
	}
}
