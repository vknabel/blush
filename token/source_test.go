package token_test

import (
	"testing"

	"github.com/vknabel/lithia/token"
)

func TestMakeSource(t *testing.T) {
	src := token.MakeSource("foo", "bar", 42)
	if src.ModuleName != "foo" {
		t.Errorf("expected %q, got %q", "foo", src.ModuleName)
	}
	if src.FileName != "bar" {
		t.Errorf("expected %q, got %q", "bar", src.FileName)
	}
	if src.Offset != 42 {
		t.Errorf("expected %d, got %d", 42, src.Offset)
	}
}
