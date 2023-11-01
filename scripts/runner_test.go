package scripts

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookPath(t *testing.T) {
	tests := []struct {
		scriptFile string
		wantErr    bool
	}{
		{"no-such-file", true},
		{"fetch-uj-records.sh", false},
		{"get-uid-map.sh", false},
		{"splunk-to-segment.sh", false},
		{"segment-mass-uploader.sh", false},
	}
	for _, tt := range tests {
		t.Run(tt.scriptFile, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			got, err := LookPath(tt.scriptFile)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err, "Must be able to find %v", tt.scriptFile)
			require.True(strings.HasSuffix(got, tt.scriptFile), "Must end with %v", tt.scriptFile)
			filNfo, err := os.Stat(got)
			require.NoError(err, "stat() call on result failed with: %v", err)
			assert.True(filNfo.Mode().IsRegular(), "result must be a regular file")
			assert.NotEqual(0, filNfo.Mode()&0444, "result must be executable")
			t.Logf("LookPath() = %v", got)
		})
	}
}
