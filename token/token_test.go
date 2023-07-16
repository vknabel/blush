package token_test

import (
	"testing"

	"github.com/vknabel/lithia/token"
)

func TestLookupIdent(t *testing.T) {
	testCases := []struct {
		input    string
		expected token.TokenType
	}{
		{"foo", token.IDENT},
		{"bar", token.IDENT},
		{"true", token.IDENT},
		{"module", token.MODULE},
		{"import", token.IMPORT},
		{"data", token.DATA},
		{"extern", token.EXTERN},
		{"func", token.FUNCTION},
		{"let", token.LET},
		{"type", token.TYPE},
		{"return", token.RETURN},
		{"if", token.IF},
		{"else", token.ELSE},
		{"for", token.FOR},
		{"_", token.BLANK},
	}

	for _, tc := range testCases {
		if tok := token.LookupIdent(tc.input); tok != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, tok)
		}
	}
}
