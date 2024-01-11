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
import json

@json.Type(json.Null)
data None
// <- ast.DeclData

@json.Inline
data Some {
// <- ast.DeclData
   value
// ^ ast.DeclField
}

@json.Inline()
enum Optional {
	None
	Some
}

extern blank // this is an extern constant
// <- ast.DeclExternFunc

extern doSomething()
// <- ast.DeclExternFunc
	
extern doSomething(argument)
// <- ast.DeclExternFunc
//     ^ ast.Identifier
//                 ^ ast.DeclParameter

extern SomeType {
// <- ast.DeclExternType
    name
//  ^ ast.DeclField
}

extern SomeEmptyType {}
// <- ast.DeclExternType
//     ^ ast.Identifier

annotation Type {
// <- ast.DeclAnnotation
	@AnyType
	value
//  ^ ast.DeclField
}

annotation ValidationRule {
	@Type(Function)
    isValid(value)
//  ^ ast.DeclField
//          ^ ast.DeclParameter
}
`

	lex := lexer.New("testmodule", "testfile.lithia", contents)
	p := parser.New(lex)
	sourceFile := p.ParseSourceFile("testfile.lithia", "testmodule", contents)
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			src := err.Source()
			contentsBeforeOffset := contents[:src.Offset]
			loc := strings.Count(contentsBeforeOffset, "\n")
			lastLineIndex := strings.LastIndex(contentsBeforeOffset, "\n")
			col := src.Offset - lastLineIndex
			relevantLine, _, _ := strings.Cut(contents[lastLineIndex+1:], "\n")

			t.Errorf("%s:%d:%d: %s\n\n  %s\n  %s^\n  %s", err.Source().FileName, loc, col, err.Summary(), relevantLine, strings.Repeat(" ", col-1), err.Details())
		}
	}

	h := syncheck.NewHarness(func(lineOffset int, line string, assert syncheck.Assertion) bool {
		var relevantChildren []ast.Node
		sourceFile.EnumerateChildNodes(func(child ast.Node) {
			tok := child.TokenLiteral()

			if tok.Source.Offset <= assert.SourceOffset && assert.SourceOffset <= tok.Source.Offset+len(tok.Literal) {
				relevantChildren = append(relevantChildren, child)
			}
		})
		for _, child := range relevantChildren {
			candidate := strings.TrimPrefix(fmt.Sprintf("%T", child), "*")
			if candidate == assert.Value {
				return !assert.Negated
			}
		}
		childTypes := make([]string, len(relevantChildren))
		for i, child := range relevantChildren {
			childTypes[i] = strings.TrimPrefix(fmt.Sprintf("%T", child), "*")
		}
		t.Errorf("no alternative found, want %q, got one of %q", assert.Value, childTypes)
		return false
	})
	err := h.Test(contents)
	if err != nil {
		t.Error(err)
	}
}
