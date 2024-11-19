package vm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/compiler"
	"github.com/vknabel/lithia/lexer"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/vm"
)

type vmTestCase struct {
	input    string
	expected any
	err      string
}

func TestBasicOperations(t *testing.T) {
	tests := []vmTestCase{
		{input: "1", expected: 1},
		{input: "1+2", expected: 3},
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "!true", expected: false},
		{input: "!false", expected: true},
		{input: "true && true", expected: true},
		{input: "true && 3", err: `unexpected type (runtime.Int "3")`},
		{input: "(if true { 2 } else { 3 })", expected: 2},
		{input: "(if 1 == 1 { 2*3 } else { 3 })", expected: 6},
		{input: "(if 1 == 0 { 2*3 } else { 3 })", expected: 3},
		{input: "(if 1 != 0 { 2*3 } else { 3 })", expected: 6},
		{input: "(if true { 2*3 } else { 3 })", expected: 6},
		{input: "(if true || false { 2*3 } else { 3 })", expected: 6},
		{input: "if true || false { 2*3 } else { 3 }", expected: 6},
		{input: `"abc"`, expected: "abc"},
		{input: "[]", expected: []any{}},
		{input: "[1, 2, 3]", expected: []any{1, 2, 3}},
		{input: "[:]", expected: map[any]any{}},
		{input: `["hello": "world", 1: 2]`, expected: map[any]any{"hello": "world", 1: 2}},
		{input: `["1": 3, 1: 2]`, expected: map[any]any{"1": 3, 1: 2}},
	}

	runVmTests(t, tests)
}
func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d.", i), func(t *testing.T) {
			program := prepareSourceFileParsing(t, tt.input)

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := vm.New(comp.Bytecode())
			err = vm.Run()
			if err != nil && tt.err == "" {
				t.Fatalf("vm error: %s", err)
			}

			if tt.err != "" {
				if err == nil || err.Error() != tt.err {
					t.Errorf("expected error %q, got %q", tt.err, err)
				}
			}
			if tt.expected != nil {
				stackElem := vm.LastPoppedStackElem()

				testExpectedValue(t, tt.expected, stackElem)
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

func testExpectedValue(t *testing.T, expected interface{}, actual runtime.RuntimeValue) {
	t.Helper()
	err := testValue(expected, actual)
	if err != nil {
		t.Error(err)
	}
}

func testValue(expected interface{}, actual runtime.RuntimeValue) error {
	switch expected := expected.(type) {
	case int:
		return testInt(int64(expected), actual)
	case bool:
		return testBool(bool(expected), actual)
	case string:
		return testString(expected, actual)
	case []any:
		return testArray([]any(expected), actual)
	case map[any]any:
		return testDict(map[any]any(expected), actual)
	default:
		return fmt.Errorf("unhandled type %T", expected)
	}
}

func testInt(expected int64, actual runtime.RuntimeValue) error {
	result, ok := actual.(runtime.Int)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if int64(result) != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result, expected)
	}

	return nil
}

func testBool(expected bool, actual runtime.RuntimeValue) error {
	result, ok := actual.(runtime.Bool)
	if !ok {
		return fmt.Errorf("object is not Bool. got=%T (%+v)", actual, actual)
	}

	if bool(result) != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result, expected)
	}

	return nil
}

func testString(expected string, actual runtime.RuntimeValue) error {
	result, ok := actual.(runtime.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if string(result) != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result, expected)
	}

	return nil
}

func testArray(expected []any, actual runtime.RuntimeValue) error {
	result, ok := actual.(runtime.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)", actual, actual)
	}

	if len(expected) != len(result) {
		return fmt.Errorf("length does not match. got=%d, want=%d", len(result), len(expected))
	}
	for i, el := range result {
		err := testValue(expected[i], el)
		if err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}
	}
	return nil
}

func testDict(expected map[any]any, actual runtime.RuntimeValue) error {
	result, ok := actual.(runtime.Dict)
	if !ok {
		return fmt.Errorf("object is not Dict. got=%T (%+v)", actual, actual)
	}

	if len(expected) != len(result) {
		return fmt.Errorf("length does not match. got=%d, want=%d", len(result), len(expected))
	}

	for key, el := range result {
		nkey, err := native(key)
		if err != nil {
			return fmt.Errorf("at index %q: %w", key, err)
		}
		err = testValue(expected[nkey], el)
		if err != nil {
			return fmt.Errorf("at index %q: %w", key, err)
		}
	}
	return nil
}

func native(val runtime.RuntimeValue) (any, error) {
	switch val := val.(type) {
	case runtime.Bool:
		return bool(val), nil
	case runtime.Int:
		return int(val), nil
	case runtime.String:
		return string(val), nil
	default:
		return nil, fmt.Errorf("cannot convert %T into native Go type, got=%q", val, val.Inspect())
	}
}
