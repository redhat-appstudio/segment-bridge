package querygen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testes in this file are simply making sure we're getting queries back
// which means we're passing in valid field names. Our query generation code
// already makes sure that the queries we generate are valid

func TestGenApplicationQuery(t *testing.T) {
	out := GenApplicationQuery("some_index")
	assert.NotEqual(t, "" , out)
}

func TestGenPipelineRunQuery(t *testing.T) {
	out := GenPipelineRunQuery("some_index")
	assert.NotEqual(t, "", out)
}
