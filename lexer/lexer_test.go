package lexer_test

import (
	"testing"

	"github.com/vknabel/lithia/lexer"
	"github.com/vknabel/lithia/token"
)

func TestLexer(t *testing.T) {
	input := `
#!/usr/bin/env lithia
module example

import tests {
	test
}
import tests.tests_t

test "any in enums matches all types", { fail =>
	enum AnyEnum {
		Int
		Any
	}

	let isCorrect = with "should be any", type AnyEnum {
		Int: { _ => False },
		Any: { _ => True }
	}
	unless isCorrect, fail "should not be the case"
}
`
	l := lexer.New("lexer", "TestLexer", input)

	expect := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.COMMENT, "!/usr/bin/env lithia"},
		{token.MODULE, "module"},
		{token.IDENT, "example"},
		{token.IMPORT, "import"},
		{token.IDENT, "tests"},
		{token.LBRACE, "{"},
		{token.IDENT, "test"},
		{token.RBRACE, "}"},
		{token.IMPORT, "import"},
		{token.IDENT, "tests"},
		{token.DOT, "."},
		{token.IDENT, "tests_t"},
		{token.IDENT, "test"},
		{token.STRING, "any in enums matches all types"},
		{token.COMMA, ","},
		{token.LBRACE, "{"},
		{token.IDENT, "fail"},
		{token.ARROW, "=>"},
		{token.ENUM, "enum"},
		{token.IDENT, "AnyEnum"},
		{token.LBRACE, "{"},
		{token.IDENT, "Int"},
		{token.IDENT, "Any"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "isCorrect"},
		{token.ASSIGN, "="},
		{token.IDENT, "with"},
		{token.STRING, "should be any"},
		{token.COMMA, ","},
		{token.TYPE, "type"},
		{token.IDENT, "AnyEnum"},
		{token.LBRACE, "{"},
		{token.IDENT, "Int"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.BLANK, "_"},
		{token.ARROW, "=>"},
		{token.IDENT, "False"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.IDENT, "Any"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.BLANK, "_"},
		{token.ARROW, "=>"},
		{token.IDENT, "True"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.IDENT, "unless"},
		{token.IDENT, "isCorrect"},
		{token.COMMA, ","},
		{token.IDENT, "fail"},
		{token.STRING, "should not be the case"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	for i, tt := range expect {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestAllTokens(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}
	}{
		{
			name:  "spaces",
			input: `    `,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EOF, ""},
			},
		},
		{
			name:  "tabs",
			input: "\t\t\t\t",
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EOF, ""},
			},
		},
		{
			name:  "newlines",
			input: "\n\n\n\n",
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EOF, ""},
			},
		},
		{
			name:  "general whitespace",
			input: "  \t\n\r\n\t  \n\t",
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EOF, ""},
			},
		},
		{
			name:  "bang",
			input: `!`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.BANG, "!"},
				{token.EOF, ""},
			},
		},
		{
			name:  "neq",
			input: `!=`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.NEQ, "!="},
				{token.EOF, ""},
			},
		},
		{
			name:  "lt",
			input: `<`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LT, "<"},
				{token.EOF, ""},
			},
		},
		{
			name:  "lte",
			input: `<=`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LTE, "<="},
				{token.EOF, ""},
			},
		},
		{
			name:  "gt",
			input: `>`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.GT, ">"},
				{token.EOF, ""},
			},
		},
		{
			name:  "gte",
			input: `>=`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.GTE, ">="},
				{token.EOF, ""},
			},
		},
		{
			name:  "plus",
			input: `+`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.PLUS, "+"},
				{token.EOF, ""},
			},
		},
		{
			name:  "minus",
			input: `-`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.MINUS, "-"},
				{token.EOF, ""},
			},
		},
		{
			name:  "asterisk",
			input: `*`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ASTERISK, "*"},
				{token.EOF, ""},
			},
		},
		{
			name:  "slash",
			input: `/`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.SLASH, "/"},
				{token.EOF, ""},
			},
		},
		{
			name:  "comma",
			input: `,`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.COMMA, ","},
				{token.EOF, ""},
			},
		},
		{
			name:  "eq",
			input: `==`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EQ, "=="},
				{token.EOF, ""},
			},
		},
		{
			name:  "assign",
			input: `=`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ASSIGN, "="},
				{token.EOF, ""},
			},
		},
		{
			name:  "fat arrow",
			input: `=>`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ARROW, "=>"},
				{token.EOF, ""},
			},
		},
		{
			name:  "thin arrow",
			input: `->`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ARROW, "->"},
				{token.EOF, ""},
			},
		},
		{
			name:  "and",
			input: `&&`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.AND, "&&"},
				{token.EOF, ""},
			},
		},
		{
			name:  "or",
			input: `||`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.OR, "||"},
				{token.EOF, ""},
			},
		},
		{
			name:  "colon",
			input: `:`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.COLON, ":"},
				{token.EOF, ""},
			},
		},
		{
			name:  "dot",
			input: `.`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.DOT, "."},
				{token.EOF, ""},
			},
		},
		{
			name:  "comma",
			input: `,`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.COMMA, ","},
				{token.EOF, ""},
			},
		},
		{
			name:  "lparen",
			input: `(`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LPAREN, "("},
				{token.EOF, ""},
			},
		},
		{
			name:  "rparen",
			input: `)`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.RPAREN, ")"},
				{token.EOF, ""},
			},
		},
		{
			name:  "lbrace",
			input: `{`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LBRACE, "{"},
				{token.EOF, ""},
			},
		},
		{
			name:  "rbrace",
			input: `}`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			name:  "lbracket",
			input: `[`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LBRACKET, "["},
				{token.EOF, ""},
			},
		},
		{
			name:  "rbracket",
			input: `]`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.RBRACKET, "]"},
				{token.EOF, ""},
			},
		},
		{
			name:  "at",
			input: `@`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.AT, "@"},
				{token.EOF, ""},
			},
		},
		{
			name:  "hash comment",
			input: `# abc def`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.COMMENT, " abc def"},
				{token.EOF, ""},
			},
		},
		{
			name:  "line comment",
			input: `// abc def`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.COMMENT, " abc def"},
				{token.EOF, ""},
			},
		},
		{
			name:  "string",
			input: `"abc"`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.STRING, "abc"},
				{token.EOF, ""},
			},
		},
		{
			name:  "EOF",
			input: ``,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EOF, ""},
			},
		},
		{
			name:  "identifier",
			input: `abc`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IDENT, "abc"},
				{token.EOF, ""},
			},
		},
		{
			name:  "identifier with underscore",
			input: `abc_def`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IDENT, "abc_def"},
				{token.EOF, ""},
			},
		},
		{
			name:  "identifier with number",
			input: `abc123`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IDENT, "abc123"},
				{token.EOF, ""},
			},
		},
		{
			name:  "identifier with number and underscore",
			input: `abc_123`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IDENT, "abc_123"},
				{token.EOF, ""},
			},
		},
		{
			name:  "identifier with camel case",
			input: `abcDef`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IDENT, "abcDef"},
				{token.EOF, ""},
			},
		},
		{
			name:  "simple integer",
			input: `42`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.INT, "42"},
				{token.EOF, ""},
			},
		},
		{
			name:  "simple float",
			input: `123.0`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.FLOAT, "123.0"},
				{token.EOF, ""},
			},
		},
		{
			name:  "integer followed by dot",
			input: `123.`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.INT, "123"},
				{token.DOT, "."},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword module",
			input: `module`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.MODULE, "module"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword import",
			input: `import`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IMPORT, "import"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword enum",
			input: `enum`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ENUM, "enum"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword data",
			input: `data`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.DATA, "data"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword extern",
			input: `extern`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.EXTERN, "extern"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword func",
			input: `func`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.FUNCTION, "func"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword let",
			input: `let`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LET, "let"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword type",
			input: `type`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.TYPE, "type"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword return",
			input: `return`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.RETURN, "return"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword if",
			input: `if`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.IF, "if"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword else",
			input: `else`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ELSE, "else"},
				{token.EOF, ""},
			},
		},
		{
			name:  "keyword for",
			input: `for`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.FOR, "for"},
				{token.EOF, ""},
			},
		},
		{
			name:  "illegal and",
			input: `&`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ILLEGAL, "&"},
				{token.EOF, ""},
			},
		},
		{
			name:  "illegal or",
			input: `|`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ILLEGAL, "|"},
				{token.EOF, ""},
			},
		},
		{
			name:  "emoji",
			input: `🦜`,
			expected: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ILLEGAL, string("🦜"[0])},
				{token.ILLEGAL, string("🦜"[1])},
				{token.ILLEGAL, string("🦜"[2])},
				{token.ILLEGAL, string("🦜"[3])},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range testCases {
		l := lexer.New("lexer", "TestAllTokens", tt.input)

		for i, expect := range tt.expected {
			tok := l.NextToken()

			if tok.Type != expect.expectedType && tok.Literal != expect.expectedLiteral {
				t.Errorf("%s[%d] - wrong token and literal. expected=%v, got=%v",
					tt.name, i, expect, tok)
				break
			}

			if tok.Type != expect.expectedType {
				t.Errorf("%s[%d] - tokentype wrong. expected=%q, got=%q",
					tt.name, i, expect.expectedType, tok.Type)
				break
			}

			if tok.Literal != expect.expectedLiteral {
				t.Errorf("%s[%d] - literal wrong. expected=%q, got=%q",
					tt.name, i, expect.expectedType, tok.Literal)
				break
			}
		}
	}
}
