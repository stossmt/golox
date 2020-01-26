package lib

import (
	"fmt"
	"strconv"
	"unicode"
)

var (
	keywords map[string]int = map[string]int{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
)

func Scan(str string, reporter *Reporter) []Token {
	s := scanner{src: []rune(str)}

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken(reporter)
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

type scanner struct {
	tokens []Token

	src     []rune
	start   int
	current int
	line    int
}

func (s *scanner) scanToken(reporter *Reporter) {
	switch r := s.advance(); r {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '=':
		s.pick('=', EQUAL_EQUAL, EQUAL)
	case '!':
		s.pick('=', BANG_EQUAL, BANG)
	case '<':
		s.pick('=', LESS_EQUAL, LESS)
	case '>':
		s.pick('=', GREATER_EQUAL, GREATER)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\t':
	case '\r':
	case '\n':
		s.line++
	case '"':
		s.string(reporter)
	default:
		if unicode.IsDigit(r) {
			s.number(reporter)
		} else if s.isAlphaNum(r) {
			s.identifier()
		} else {
			msg := fmt.Sprintf("unexpected character: %v\n\n", string(r))
			reporter.Report(s.line, "", msg)
		}
	}
}

func (s *scanner) advance() rune {
	s.current++
	return s.src[s.current-1]
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.src[s.current]
}

func (s *scanner) peekNext() rune {
	if s.current+1 >= len(s.src) {
		return '\000'
	}
	return s.src[s.current+1]
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *scanner) isAlphaNum(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

func (s *scanner) addToken(tokType int) {
	s.addLiteral(tokType, nil)
}

func (s *scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	} else if s.src[s.current] != r {
		return false
	} else {
		return true
	}
}

func (s *scanner) pick(r rune, success int, failure int) {
	if s.match(r) {
		s.current++
		s.addToken(success)
	} else {
		s.addToken(failure)
	}
}

func (s *scanner) string(reporter *Reporter) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		reporter.Report(s.line, "", "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	literal := string(s.src[s.start+1 : s.current-1])
	s.addLiteral(STRING, literal)
}

func (s *scanner) number(reporter *Reporter) {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// Look for a float.
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	str := string(s.src[s.start:s.current])
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid number literal: %v", str)
		reporter.Report(s.line, "", msg)
	}
	s.addLiteral(NUMBER, num)
}

func (s *scanner) identifier() {
	for s.isAlphaNum(s.peek()) {
		s.advance()
	}

	str := string(s.src[s.start:s.current])
	if keyword_type, ok := keywords[str]; ok {
		s.addToken(keyword_type)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *scanner) addLiteral(tokType int, literal interface{}) {
	lexeme := string(s.src[s.start:s.current])
	t := Token{tokType, lexeme, literal, s.line}
	s.tokens = append(s.tokens, t)
}
