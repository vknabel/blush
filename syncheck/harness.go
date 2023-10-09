package syncheck

import (
	"fmt"
	"strings"
)

type SyntaxMatcher func(line string, assert Assertion) bool

type Harness struct {
	match SyntaxMatcher
}

func NewHarness(matcher SyntaxMatcher) Harness {
	return Harness{matcher}
}

func (h *Harness) Test(doc string) error {
	asserts := ParseAssertions(doc)
	lines := strings.Split(doc, "\n")
	var failures []Assertion
	for _, a := range asserts {
		matched := h.match(lines[a.Line-1], a)
		if matched == a.Negated {
			failures = append(failures, a)
		}
	}

	if len(failures) == 0 {
		return nil
	}
	var out strings.Builder
	out.WriteString(fmt.Sprintf("failed assertions %d:\n", len(failures)))

	for i, a := range failures {
		out.WriteString("\n")
		out.WriteString(fmt.Sprintf("%d. error:\n", i+1))
		out.WriteString("\t" + lines[a.Line-1] + "\n")
		out.WriteString("\t" + lines[a.SourceLine-1] + "\n")
	}
	return fmt.Errorf(out.String())
}
