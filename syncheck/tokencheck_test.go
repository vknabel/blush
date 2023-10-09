package syncheck_test

import (
	"testing"

	"github.com/vknabel/lithia/syncheck"
)

func TestParsingAssertions(t *testing.T) {
	content := `
	data xyz
	`
	asserts := syncheck.ParseAssertions(content)
	if len(asserts) != 0 {
		t.Errorf("no asserts expected, got: %v", asserts)
	}
	content = "data xyz {\n" +
		"//   ^ label0\n" +
		"// <- label1\n" +
		"}\n" +
		"^ label2"

	asserts = syncheck.ParseAssertions(content)
	if len(asserts) != 3 {
		t.Errorf("three assertions expected, got: %v", asserts)
	}
	if asserts[0].Line != 1 || asserts[0].Column != 6 {
		t.Errorf("expected matcher [0] on line 1 column 6, got: %+v", asserts[0])
	}
	if asserts[1].Line != 1 || asserts[1].Column != 1 {
		t.Errorf("expected matcher [1] on line 1 column 1, got: %+v", asserts[1])
	}
	if asserts[2].Line != 4 || asserts[1].Column != 1 {
		t.Errorf("expected matcher [2] on line 4 column 1, got: %+v", asserts[2])
	}
}
