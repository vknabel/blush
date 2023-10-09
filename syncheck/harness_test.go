package syncheck_test

import (
	"strings"
	"testing"

	"github.com/vknabel/lithia/syncheck"
)

func TestHarness(t *testing.T) {
	h := syncheck.NewHarness(func(line string, assert syncheck.Assertion) bool {
		identifier := line[assert.Column-1:]
		identifier, _, _ = strings.Cut(identifier, " ")
		return strings.ToUpper(identifier) == assert.Value
	})
	err := h.Test(`
data value {
//   ^ VALUE
}
var
// <- VAR

negated
^ !true
`)
	if err != nil {
		t.Error(err)
	}
}
