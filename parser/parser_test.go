package parser_test

import (
	"testing"

	"github.com/vknabel/lithia/syncheck"
)

func TestParseSourceFile(t *testing.T) {
	h := syncheck.NewHarness(func(line string, assert syncheck.Assertion) bool {
		return false
	})

	err := h.Test(`
	data None {}
	// <- ast.DeclData
	`)
	if err != nil {
		t.Error(err)
	}
}
