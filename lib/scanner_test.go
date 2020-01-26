package lib

import "testing"

func TestSingleSuccess(t *testing.T) {
	input := "(){},.-+;*/"
	want := []Token{
		Token{LEFT_PAREN, "(", nil, 0},
		Token{RIGHT_PAREN, ")", nil, 0},
		Token{LEFT_BRACE, "{", nil, 0},
		Token{RIGHT_BRACE, "}", nil, 0},
		Token{COMMA, ",", nil, 0},
		Token{DOT, ".", nil, 0},
		Token{MINUS, "-", nil, 0},
		Token{PLUS, "+", nil, 0},
		Token{SEMICOLON, ";", nil, 0},
		Token{STAR, "*", nil, 0},
		Token{SLASH, "/", nil, 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestSingleOrDoubleSuccess(t *testing.T) {
	str_tok_map := map[string]Token{
		"=":  Token{EQUAL, "=", nil, 0},
		"==": Token{EQUAL_EQUAL, "==", nil, 0},
		"!":  Token{BANG, "!", nil, 0},
		"!=": Token{BANG_EQUAL, "!=", nil, 0},
		"<":  Token{LESS, "<", nil, 0},
		"<=": Token{LESS_EQUAL, "<=", nil, 0},
		">":  Token{GREATER, ">", nil, 0},
		">=": Token{GREATER_EQUAL, ">=", nil, 0},
	}

	end := Token{EOF, "", nil, 0}

	for input, want := range str_tok_map {
		reporter := NewReporter()
		tokens := Scan(input, reporter)
		if reporter.HadErr {
			t.Errorf("Scanner reported an error")
		}
		if tokens[0] != want {
			t.Errorf("Wrong token; got: %v, want: %v", tokens[0], want)
		}
		if tokens[1] != end {
			t.Errorf("Wrong token; got: %v, want: %v", tokens[1], end)
		}
	}
}

func TestOneCommentSuccess(t *testing.T) {
	input := "// A comment should be skipped."
	want := Token{EOF, "", nil, 0}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if tokens[0] != want {
		t.Errorf("Comment was not discarded; got %v, want %v", tokens[0], want)
	}
}

func TestMultipleCommentsSuccess(t *testing.T) {
	input := "// A comment\n// Should preserve line numbers."
	want := Token{EOF, "", nil, 1}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if tokens[0] != want {
		t.Errorf("Comment was not discarded; got %v, want %v", tokens[0], want)
	}
}

func TestWhiteSpaceConsumedSuccess(t *testing.T) {
	input := "( ( ) // A comment. \n ) \r  \t"
	want := []Token{
		Token{LEFT_PAREN, "(", nil, 0},
		Token{LEFT_PAREN, "(", nil, 0},
		Token{RIGHT_PAREN, ")", nil, 0},
		Token{RIGHT_PAREN, ")", nil, 1},
		Token{EOF, "", nil, 1},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestStringLiteral(t *testing.T) {
	input := "\"Hèllo Wôrld!\""
	want := []Token{
		Token{STRING, "\"Hèllo Wôrld!\"", "Hèllo Wôrld!", 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestUnterminatedStringFailure(t *testing.T) {
	input := "\"Hèllo Wôrld!"

	reporter := NewReporter()
	Scan(input, reporter)

	if !reporter.HadErr {
		t.Errorf("Scanner did not report an unterminated string literal")
	}
}

func TestNumberArithemetic(t *testing.T) {
	input := "3 * 4 / 9.0 - 0"
	want := []Token{
		Token{NUMBER, "3", float64(3), 0},
		Token{STAR, "*", nil, 0},
		Token{NUMBER, "4", float64(4), 0},
		Token{SLASH, "/", nil, 0},
		Token{NUMBER, "9.0", float64(9), 0},
		Token{MINUS, "-", nil, 0},
		Token{NUMBER, "0", float64(0), 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestKeywords(t *testing.T) {
	input := "and class else false for fun if nil or print " +
		"return super this true var while"
	want := []Token{
		Token{AND, "and", nil, 0},
		Token{CLASS, "class", nil, 0},
		Token{ELSE, "else", nil, 0},
		Token{FALSE, "false", nil, 0},
		Token{FOR, "for", nil, 0},
		Token{FUN, "fun", nil, 0},
		Token{IF, "if", nil, 0},
		Token{NIL, "nil", nil, 0},
		Token{OR, "or", nil, 0},
		Token{PRINT, "print", nil, 0},
		Token{RETURN, "return", nil, 0},
		Token{SUPER, "super", nil, 0},
		Token{THIS, "this", nil, 0},
		Token{TRUE, "true", nil, 0},
		Token{VAR, "var", nil, 0},
		Token{WHILE, "while", nil, 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestIdentifiers(t *testing.T) {
	input := "name1 name2"
	want := []Token{
		Token{IDENTIFIER, "name1", nil, 0},
		Token{IDENTIFIER, "name2", nil, 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestIdentifierAndKeyword(t *testing.T) {
	input := "or orange or"
	want := []Token{
		Token{OR, "or", nil, 0},
		Token{IDENTIFIER, "orange", nil, 0},
		Token{OR, "or", nil, 0},
		Token{EOF, "", nil, 0},
	}

	reporter := NewReporter()
	tokens := Scan(input, reporter)

	if reporter.HadErr {
		t.Errorf("Scanner reported an error")
	}
	if !isEqual(want, tokens) {
		t.Errorf("Incorrect token sequence; got %v, want %v", tokens, want)
	}
}

func TestUnrecognizedFailure(t *testing.T) {
	reporter := NewReporter()
	Scan("^", reporter)
	if !reporter.HadErr {
		t.Errorf("Wanted error on unrecognized token")
	}
}

func isEqual(want []Token, got []Token) bool {
	if len(want) != len(got) {
		return false
	}
	for i := range want {
		if want[i] != got[i] {
			return false
		}
	}
	return true
}
