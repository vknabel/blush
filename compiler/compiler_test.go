package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vknabel/blush/ast"
	"github.com/vknabel/blush/compiler"
	"github.com/vknabel/blush/lexer"
	code "github.com/vknabel/blush/op"
	"github.com/vknabel/blush/parser"
	"github.com/vknabel/blush/registry/staticmodule"
	"github.com/vknabel/blush/runtime"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestUnaryOperators(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "!true",
			expectedConstants: nil,
			expectedInstructions: []code.Instructions{
				code.Make(code.ConstTrue),
				code.Make(code.Invert),
				code.Make(code.Pop),
			},
		},
		{
			input:             "-3",
			expectedConstants: []any{3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Negate),
				code.Make(code.Pop),
			},
		},
		{
			input:             "+42",
			expectedConstants: []any{42},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBinaryOperators(t *testing.T) {
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
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// -
				code.Make(code.Sub),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// *
				code.Make(code.Mul),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 / 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// /
				code.Make(code.Div),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// ==
				code.Make(code.Equal),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// !=
				code.Make(code.NotEqual),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// >
				code.Make(code.GreaterThan),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// <
				code.Make(code.LessThan),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// >=
				code.Make(code.GreaterThanOrEqual),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 <= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// <=
				code.Make(code.LessThanOrEqual),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 % 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				// 1
				code.Make(code.Const, 0),
				// 2
				code.Make(code.Const, 1),
				// %
				code.Make(code.Mod),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "true && false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				// left
				code.Make(code.ConstTrue),
				// when false do not exectue right
				code.Make(code.JumpFalse, 11),
				// right
				code.Make(code.ConstFalse),
				code.Make(code.AssertType, int(runtime.Bool(true).TypeConstantId())),
				// result is right
				code.Make(code.Jump, 12),
				// put false back up
				code.Make(code.ConstFalse),
				// drop expr
				code.Make(code.Pop),
			},
		},
		{
			input:             "true || false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				// left
				code.Make(code.ConstTrue),
				// when true do not exectue right
				code.Make(code.JumpTrue, 11),
				// right
				code.Make(code.ConstFalse),
				code.Make(code.AssertType, int(runtime.Bool(true).TypeConstantId())),
				// result is right
				code.Make(code.Jump, 12),
				// put true back up
				code.Make(code.ConstTrue),
				// drop expr
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TesEQtIfStmtsArithmetic(t *testing.T) {
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

func TestArrayExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []any{0},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Array),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[42, 1337]",
			expectedConstants: []any{42, 1337, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Const, 1),
				code.Make(code.Const, 2),
				code.Make(code.Array),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[42 + 1337]",
			expectedConstants: []any{42, 1337, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Const, 1),
				code.Make(code.Add),
				code.Make(code.Const, 2),
				code.Make(code.Array),
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
	l, err := lexer.New(staticmodule.NewSourceString("testing:///test/test.blush", input))
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewSourceParser(l, nil, "test.blush")
	srcFile := p.ParseSourceFile()
	checkParserErrors(t, p, input)
	return srcFile
}

func checkParserErrors(t *testing.T, p *parser.Parser, contents string) {
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			src := err.Token.Source

			if src == nil {
				t.Errorf("<no source>: %q\n  %s", err.Token.Literal, err.Details)
				continue
			}
			contentsBeforeOffset := contents[:src.Offset]
			loc := strings.Count(contentsBeforeOffset, "\n")
			lastLineIndex := strings.LastIndex(contentsBeforeOffset, "\n")
			col := src.Offset - lastLineIndex
			relevantLine, _, _ := strings.Cut(contents[lastLineIndex+1:], "\n")

			t.Errorf("%s:%d:%d: %s\n\n  %s\n  %s^\n  %s", err.Token.Source.File, loc, col, err.Summary, relevantLine, strings.Repeat(" ", col-1), err.Details)
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
		case bool:
			got, ok := actual[i].(runtime.Bool)
			if !ok || want != bool(got) {
				return fmt.Errorf("wrong constant at %d.\nwant=%t\ngot=%q", i, want, got.Inspect())
			}
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
