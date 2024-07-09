package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/compiler"
	"github.com/vknabel/lithia/lexer"
	code "github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/runtime"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// +
				code.Make(code.Add),
				// drop expr
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIfStmtsArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "if 1 { 2 } else { 3 }",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 13),
				code.Make(code.Const, 1),
				code.Make(code.Pop),
				code.Make(code.Jump, 17),
				code.Make(code.Const, 2),
				code.Make(code.Pop),
			},
		},
		{
			input:             "if 1 { 2 }",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 10),
				code.Make(code.Const, 1),
				code.Make(code.Pop),
			},
		},
		{
			input:             "if 0 { 1 } else if 2 { 3 } else { 4 }",
			expectedConstants: []interface{}{0, 1, 2, 3, 4},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 13),
				code.Make(code.Const, 1),
				code.Make(code.Pop),
				code.Make(code.Jump, 30),
				code.Make(code.Const, 2),
				code.Make(code.JumpFalse, 26),
				code.Make(code.Const, 3),
				code.Make(code.Pop),
				code.Make(code.Jump, 30),
				code.Make(code.Const, 4),
				code.Make(code.Pop),
			},
		},
		{
			input:             "if 0 { 1 } else if 2 { 3 }",
			expectedConstants: []interface{}{0, 1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 13),
				code.Make(code.Const, 1),
				code.Make(code.Pop),
				code.Make(code.Jump, 23),
				code.Make(code.Const, 2),
				code.Make(code.JumpFalse, 23),
				code.Make(code.Const, 3),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIfExpressionsArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "(if 1 { 2 } else { 3 })",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 12),
				code.Make(code.Const, 1),
				code.Make(code.Jump, 15),
				code.Make(code.Const, 2),
				code.Make(code.Pop),
			},
		},
		{
			input:             "(if 0 { 1 } else if 2 { 3 } else { 4 })",
			expectedConstants: []interface{}{0, 1, 2, 3, 4},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.JumpFalse, 12),
				code.Make(code.Const, 1),
				code.Make(code.Jump, 27),
				code.Make(code.Const, 2),
				code.Make(code.JumpFalse, 24),
				code.Make(code.Const, 3),
				code.Make(code.Jump, 27),
				code.Make(code.Const, 4),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d.", i), func(t *testing.T) {
			program := prepareSourceFileParsing(t, tt.input)

			compiler := compiler.New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			bytecode := compiler.Bytecode()

			err = testInstructions(t, tt.expectedInstructions, bytecode.Instructions)
			if err != nil {
				t.Fatalf("testInstructions failed: %s", err)
			}

			err = testConstants(t, tt.expectedConstants, bytecode.Constants)
			if err != nil {
				t.Fatalf("testConstants failed: %s", err)
			}
		})
	}
}

func prepareSourceFileParsing(t *testing.T, input string) *ast.SourceFile {
	l := lexer.New("testing", "test.lithia", input)
	p := parser.New(l)
	srcFile := p.ParseSourceFile("test.lithia", "testing", nil, input)
	checkParserErrors(t, p, input)
	return srcFile
}

func checkParserErrors(t *testing.T, p *parser.Parser, contents string) {
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			src := err.Token.Source
			contentsBeforeOffset := contents[:src.Offset]
			loc := strings.Count(contentsBeforeOffset, "\n")
			lastLineIndex := strings.LastIndex(contentsBeforeOffset, "\n")
			col := src.Offset - lastLineIndex
			relevantLine, _, _ := strings.Cut(contents[lastLineIndex+1:], "\n")

			t.Errorf("%s:%d:%d: %s\n\n  %s\n  %s^\n  %s", err.Token.Source.FileName, loc, col, err.Summary, relevantLine, strings.Repeat(" ", col-1), err.Details)
		}
		t.FailNow()
	}
}

func testInstructions(
	t *testing.T,
	expected []code.Instructions,
	actual code.Instructions,
) error {
	t.Helper()
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testConstants(
	t *testing.T,
	expected []any,
	actual []runtime.RuntimeValue,
) error {
	t.Helper()

	if len(actual) != len(expected) {
		return fmt.Errorf("wrong amount of constants.\nwant=%q\ngot=%q", expected, actual)
	}

	for i, cons := range expected {
		switch want := cons.(type) {
		case int:
			got, ok := actual[i].(runtime.Int)
			if !ok || want != int(got) {
				return fmt.Errorf("wrong constant at %d.\nwant=%d\ngot=%q", i, want, got)
			}
		default:
			got := actual[i]
			return fmt.Errorf("unhandled wanted type %T of value at %d.\nwant=%q\ngot=%q", i, want, want, got)
		}
	}
	return nil
}
