package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/lexer"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/syncheck"
)

func TestParseSourceFile(t *testing.T) {
	contents := `
	data None {}
	// <- ast.DeclData
	`

	lex := lexer.New("testmodule", "testfile.lithia", contents)
	p := parser.New(lex)
	sourceFile := p.ParseSourceFile("testfile.lithia", "testmodule", contents)
	if len(p.Errors()) > 0 {
		t.Error(p.Errors())
	}

	h := syncheck.NewHarness(func(lineOffset int, line string, assert syncheck.Assertion) bool {
		var relevantChildren []ast.Node
		sourceFile.EnumerateChildNodes(func(child ast.Node) {
			tok := child.TokenLiteral()
			fmt.Printf("l%d o%d l%s:child: %T, %+v", lineOffset, assert.SourceOffset, line, child, child.TokenLiteral().Source)
			if tok.Source.Offset <= assert.SourceOffset && assert.SourceOffset <= tok.Source.Offset+len(tok.Literal) {
				relevantChildren = append(relevantChildren, child)
			}
		})
		t.Logf("relevant children: %d", len(relevantChildren))
		for _, child := range relevantChildren {
			if strings.TrimPrefix(fmt.Sprintf("%T", child), "*") == assert.Value {
				return !assert.Negated
			}
		}
		return false
	})
	err := h.Test(contents)
	if err != nil {
		t.Error(err)
	}
}
