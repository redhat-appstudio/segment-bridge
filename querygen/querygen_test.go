package querygen

import (
	"testing"
)

func TestGenQuery(t *testing.T) {
	expected := "Hi!"
	out := GenQuery()
	if out != expected {
		t.Errorf("GenQuery() == %q, expected %q", out, expected)
	}
}
