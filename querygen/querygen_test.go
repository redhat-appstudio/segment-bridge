package querygen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests in this file are simply making sure we're getting queries back
// which means we're passing in valid field names. Our query generation code
// already makes sure that the queries we generate are valid

func TestGenApplicationQuery(t *testing.T) {
	out := GenApplicationQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenBuildPipelineRunCreatedQuery(t *testing.T) {
	out := GenBuildPipelineRunCreatedQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenBuildPipelineRunStartedQuery(t *testing.T) {
	out := GenBuildPipelineRunStartedQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenClairScanCompletedQuery(t *testing.T) {
	out := GenClairScanCompletedQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenBuildPipelineRunCompletedQuery(t *testing.T) {
	out := GenBuildPipelineRunCompletedQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenReleaseCompletedQuery(t *testing.T) {
	out := GenReleaseCompletedQuery("some_index")
	assert.NotEqual(t, "", out)
}

func TestGenPullRequestCreatedQuery(t *testing.T) {
	out := GenPullRequestCreatedQuery("some_index")
	assert.NotEqual(t, "", out)
}
