package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"testing/quick"

	"github.com/lithammer/dedent"
	"github.com/redhat-appstudio/segment-bridge.git/scripts"
	"github.com/redhat-appstudio/segment-bridge.git/stats"
	"github.com/redhat-appstudio/segment-bridge.git/webfixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Our uploader script shouldn't really care about the contents of the records
// we send so we can get away with emulating them as simple string maps
type testRecord map[string]string

type testCase struct {
	name         string
	data         []testRecord
	dataJsonSize int
	maxBatchSize int
	shouldSplit  bool
}

// Structure of a Segment batch record - for decoding JSON
type segmentBatch struct {
	Batch []map[string]any
}

func TestUploader(t *testing.T) {
	script, err := scripts.LookPath("segment-mass-uploader.sh")
	if err != nil {
		t.Fatalf("Failed to find script to test: %v", err)
	}
	testCases := mkTestCases(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			reqs := webfixture.TraceRequestsFrom(func(url string, _ *http.Client) {
				t.Setenv("SEGMENT_BATCH_API", url)
				t.Setenv("SEGMENT_BATCH_DATA_SIZE", fmt.Sprintf("%d", tt.maxBatchSize))

				withScriptStdin(t, script, func(stdin io.WriteCloser) {
					err := streamAsJsonLines(stdin, tt.data)
					require.NoError(t, err)
				})
			})

			if tt.shouldSplit {
				assert.Greater(t, len(reqs), 1, "Data should be split to batches")
			}

			requestRecords := 0
			for _, request := range reqs {
				assert.Equal(t, "POST", request.Method, "HTTP method must be POST")

				var reqData segmentBatch
				requireJSONDecode(
					t, request.Body, &reqData,
					"Failed to decode sent request JSON. "+
						"Perhaps not sent with the right structure?",
				)

				requestRecords += len(reqData.Batch)

				assert.LessOrEqual(t, len(request.Body), tt.maxBatchSize)
			}
			assert.Equal(t, len(tt.data), requestRecords, "Wrong number of records sent")
		})
	}
}

func streamAsJsonLines(stream io.Writer, data []testRecord) error {
	enc := json.NewEncoder(stream)
	for _, record := range data {
		if err := enc.Encode(record); err != nil {
			return err
		}
		if _, err := io.WriteString(stream, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func requireJSONDecode(t *testing.T, data string, v any, msgAndArgs ...any) {
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(v)
	require.NoError(t, err, msgAndArgs...)
	require.False(t, decoder.More())
}

func withScriptStdin(t *testing.T, script string, tFunc func(io.WriteCloser)) {
	cmd := exec.Command(script)
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to connect to script STDIN")
	require.NoError(t, cmd.Start(), "Failed to start script")
	defer func() {
		stdin.Close()
		require.NoError(t, cmd.Wait(), "Script did not finish cleanly")
	}()
	tFunc(stdin)
}

// Use the quick module to generate random test data
func mkTestCases(t *testing.T) (cases []testCase) {
	var records, bytes stats.Series[int]
	err := quick.Check(func(data []testRecord) bool {
		jsonData, err := json.Marshal(data)
		if !assert.NoError(t, err, "Failed to convert test sample to JSON") {
			return true
		}
		cases = append(cases, testCase{
			name:         fmt.Sprintf("records-%d-bytes-%d", len(data), len(jsonData)),
			data:         data,
			dataJsonSize: len(jsonData),
		})
		records.Add(len(data))
		bytes.Add(len(jsonData))
		return true
	}, nil)
	require.NoError(t, err, "Failed to generate test data")
	t.Logf(dedent.Dedent(`
		%d test cases
		Records: %5d
		Bytes:   %5d
	`), records.Len(), records, bytes)
	// Well set the maximum batch size to 45% of the maximum test data chunk
	// so that some samples will get split into 2 and 3 batches.
	maxBatchSize := bytes.Max() * 45 / 100
	for i := range cases {
		cases[i].maxBatchSize = maxBatchSize
		cases[i].shouldSplit = (cases[i].dataJsonSize > maxBatchSize)
	}
	return
}
