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
	label                string
	input                string
	expectedConstants    []interface{}
	expectedGlobals      [][]code.Instructions
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

func TestCharLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			label:             "simple char",
			input:             "'a'",
			expectedConstants: []any{'a'},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "escaped newline",
			input:             "'\\n'",
			expectedConstants: []any{'\n'},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "escaped quote",
			input:             "'\\''",
			expectedConstants: []any{'\''},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "escaped backslash",
			input:             "'\\\\'",
			expectedConstants: []any{'\\'},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestNumberLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			label:             "decimal integer",
			input:             "42",
			expectedConstants: []any{42},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "hexadecimal integer",
			input:             "0xFFF",
			expectedConstants: []any{4095}, // 0xFFF = 4095
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "octal integer",
			input:             "0777",
			expectedConstants: []any{511}, // 0777 = 511
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "binary integer",
			input:             "0b101010",
			expectedConstants: []any{42}, // 0b101010 = 42
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "float literal",
			input:             "3.14",
			expectedConstants: []any{3.14},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Pop),
			},
		},
		{
			label:             "scientific notation float",
			input:             "2e10",
			expectedConstants: []any{2e10},
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
		{
			label:             "array index",
			input:             "[1, 2, 3][1]",
			expectedConstants: []any{1, 2, 3, 3, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Const, 1),
				code.Make(code.Const, 2),
				code.Make(code.Const, 3),
				code.Make(code.Array),
				code.Make(code.Const, 4),
				code.Make(code.GetIndex),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestDictExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[:]",
			expectedConstants: []any{0},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Dict),
				code.Make(code.Pop),
			},
		},
		{
			label:             "dict with two key-value pairs",
			input:             "[1: 2, 3: 4]",
			expectedConstants: []any{1, 2, 3, 4, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Const, 1),
				code.Make(code.Const, 2),
				code.Make(code.Const, 3),
				code.Make(code.Const, 4),
				code.Make(code.Dict),
				code.Make(code.Pop),
			},
		},
		{
			label:             "dict with expressions",
			input:             "[1 + 1: 2 * 2]",
			expectedConstants: []any{1, 1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Const, 1),
				code.Make(code.Add),
				code.Make(code.Const, 2),
				code.Make(code.Const, 3),
				code.Make(code.Mul),
				code.Make(code.Const, 4),
				code.Make(code.Dict),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestDeclFunction(t *testing.T) {
	tests := []compilerTestCase{
		{
			label: "function with value return",
			input: "func example() { return 42 }",
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.Const, 1),
						code.Make(code.Return),
					},
				},
				42,
			},
			expectedInstructions: []code.Instructions{},
		},
		{
			label: "function call with value return",
			input: "func example() { return 42 }\nexample()",
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.Const, 1),
						code.Make(code.Return),
					},
				},
				42,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Call, 0),
				code.Make(code.Pop),
			},
		},
		{
			label: "function with blank return",
			input: "func example() { return }",
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.ConstNull),
						code.Make(code.Return),
					},
				},
			},
			expectedInstructions: []code.Instructions{},
		},
		{
			label: "function call with blank return",
			input: "func example() { return }\nexample()",
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.ConstNull),
						code.Make(code.Return),
					},
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 0),
				code.Make(code.Call, 0),
				code.Make(code.Pop),
			},
		},
		{
			label: "function can access global variables",
			input: `
			let x = 42
			func example() {
				return x
			}
			`,
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.GetGlobal, 0),
						code.Make(code.Return),
					},
				},
				42,
			},
			expectedGlobals: [][]code.Instructions{
				{code.Make(code.Const, 1)},
			},
			expectedInstructions: []code.Instructions{},
		},
		{
			label: "function can access global variables declared after usage",
			input: `
			func example() {
				return x
			}
			let x = 42
			`,
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.GetGlobal, 0),
						code.Make(code.Return),
					},
				},
				42,
			},
			expectedGlobals: [][]code.Instructions{
				{code.Make(code.Const, 1)},
			},
			expectedInstructions: []code.Instructions{},
		},
		{
			label: "function with local variable",
			input: `
			func example() {
				let x = 42
				return x+x
			}
			`,
			expectedConstants: []any{
				compiledFunction{
					name:   "example",
					params: 0,
					ins: []code.Instructions{
						code.Make(code.Const, 1),
						code.Make(code.SetLocal, 0),
						code.Make(code.GetLocal, 0),
						code.Make(code.GetLocal, 0),
						code.Make(code.Add),
						code.Make(code.Return),
					},
				},
				42,
			},
			expectedInstructions: []code.Instructions{},
		},
	}

	runCompilerTests(t, tests)
}

func TestVariables(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "let a = 42\na",
			expectedConstants: []any{
				42,
			},
			expectedGlobals: [][]code.Instructions{
				{code.Make(code.Const, 0)},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.GetGlobal, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestDeclData(t *testing.T) {
	tests := []compilerTestCase{
		{
			label: "empty data declaration",
			input: `data Example`,
			expectedConstants: []any{
				compiledDataType{
					name:   "Example",
					fields: []compiledField{},
				},
			},
		},
		{
			label: "data declaration",
			input: `data Example { field }`,
			expectedConstants: []any{
				compiledDataType{
					name:   "Example",
					fields: []compiledField{{name: "field"}},
				},
			},
		},
		{
			label: "data declaration and call",
			input: `
				data Person {
					name
				}
				Person("Max").name
				`,
			expectedConstants: []any{
				compiledDataType{
					name: "Person",
					fields: []compiledField{
						{name: "name"},
					},
				},
				"Max",
				"name",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 1),
				code.Make(code.Const, 0),
				code.Make(code.Call, 1),
				code.Make(code.GetField, 2),
				code.Make(code.Pop),
			},
		},
		{
			label: "data declaration and call with two fields",
			input: `
				data Person {
					name
					age
				}
				Person("Max", 42).name
				`,
			expectedConstants: []any{
				compiledDataType{
					name: "Person",
					fields: []compiledField{
						{name: "name"},
						{name: "age"},
					},
				},
				"Max",
				42,
				"name",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Const, 1),
				code.Make(code.Const, 2),
				code.Make(code.Const, 0),
				code.Make(code.Call, 2),
				code.Make(code.GetField, 3),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d. %s", i, tt.label), func(t *testing.T) {
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

			err = testGlobals(t, tt.expectedGlobals, bytecode.Globals)
			if err != nil {
				t.Fatalf("testGlobals failed: %s", err)
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
		case float64:
			got, ok := actual[i].(runtime.Float)
			if !ok || want != float64(got) {
				return fmt.Errorf("wrong constant at %d.\nwant=%f\ngot=%f", i, want, float64(got))
			}
		case rune:
			got, ok := actual[i].(runtime.Char)
			if !ok || want != rune(got) {
				return fmt.Errorf("wrong constant at %d.\nwant=%q\ngot=%q", i, string(want), got)
			}
		case string:
			got, ok := actual[i].(runtime.String)
			if !ok || want != string(got) {
				return fmt.Errorf("wrong constant at %d.\nwant=%q\ngot=%q", i, want, got)
			}
		case compiledFunction:
			got, ok := actual[i].(*runtime.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d is not a function: %T", i, actual[i])
			}

			if got.Symbol.Name != want.name {
				return fmt.Errorf("wrong function name at %d.\nwant=%q\ngot=%q", i, want.name, got.Symbol.Name)
			}

			if got.Arity() != want.params {
				return fmt.Errorf("wrong function params at %d.\nwant=%d\ngot=%d", i, want.params, got.Arity())
			}

			err := testInstructions(t, want.ins, got.Instructions)
			if err != nil {
				return fmt.Errorf("wrong function instructions at %d: %s", i, err)
			}

		case compiledDataType:
			got, ok := actual[i].(*runtime.DataType)
			if !ok {
				return fmt.Errorf("constant %d is not a data type: %T", i, actual[i])
			}

			if got.Symbol.Name != want.name {
				return fmt.Errorf("wrong data type name at %d.\nwant=%q\ngot=%q", i, want.name, got.Symbol.Name)
			}

			if len(got.FieldSymbols) != len(want.fields) {
				return fmt.Errorf("wrong amount of fields at %d.\nwant=%d\ngot=%d", i, len(want.fields), len(got.FieldSymbols))
			}

			for j, field := range want.fields {
				if got.FieldSymbols[j].Name != field.name {
					return fmt.Errorf("wrong field name at %d.%d.\nwant=%q\ngot=%q", i, j, field.name, got.FieldSymbols[j].Name)
				}
			}

		default:
			got := actual[i]
			return fmt.Errorf("unhandled wanted type %T of value at %d.\nwant=%q\ngot=%q", i, want, want, got)
		}
	}
	return nil
}

func testGlobals(t *testing.T,
	expected [][]code.Instructions,
	actual []*compiler.CompilationScope,
) error {
	t.Helper()

	if len(actual) != len(expected) {
		return fmt.Errorf("wrong amount of globals.\nwant=%d\ngot=%d", len(expected), len(actual))
	}

	for i, ins := range expected {
		err := testInstructions(t, ins, actual[i].Instructions)
		if err != nil {
			return fmt.Errorf("wrong global instructions at %d: %s", i, err)
		}
	}

	return nil
}

type compiledFunction struct {
	name   string
	params int
	ins    []code.Instructions
}
type compiledDataType struct {
	name   string
	fields []compiledField
}

type compiledField struct {
	name string
}
