package latex

import (
	"strconv"
)

// TokenType represents the type of a token
type TokenType int

const (
	NUMBER   TokenType = iota // 数値リテラル
	PLUS                      // +
	MINUS                     // -
	MULTIPLY                  // *
	DIVIDE                    // /
	LPAREN                    // (
	RPAREN                    // )
	VARIABLE                  // 変数（a-z, A-Z）
	CARET                     // ^
	LBRACE                    // {
	RBRACE                    // }
	LBRACKET                  // [
	RBRACKET                  // ]
	COMMAND                   // \sqrt, etc.
	EQUAL                     // =
	EOF                       // 入力終端
	ILLEGAL                   // 不正なトークン
)

// Token represents a lexical token
type Token struct {
	Type    TokenType // トークンの種類
	Literal string    // 元のテキスト
	Value   float64   // 数値の場合の値
	Pos     int       // 入力文字列内の位置
}

// Lexer performs lexical analysis on LaTeX input
type Lexer struct {
	input        string // 入力文字列
	position     int    // 現在の位置
	readPosition int    // 次に読む位置
	ch           byte   // 現在の文字
}

// NewLexer creates a new Lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NUL character indicates EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips over whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter checks if a character is a letter (a-z, A-Z)
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// readCommand reads a LaTeX command starting with backslash
func (l *Lexer) readCommand() string {
	startPos := l.position
	l.readChar() // skip '\'

	// Read command name (letters only)
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPos+1 : l.position]
}

// readNumber reads a number (integer or decimal)
func (l *Lexer) readNumber() string {
	startPos := l.position

	// Read digits before decimal point
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'

		// Read digits after decimal point
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[startPos:l.position]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Pos = l.position

	switch l.ch {
	case '+':
		tok = Token{Type: PLUS, Literal: "+", Pos: l.position}
		l.readChar()
	case '-':
		tok = Token{Type: MINUS, Literal: "-", Pos: l.position}
		l.readChar()
	case '*':
		tok = Token{Type: MULTIPLY, Literal: "*", Pos: l.position}
		l.readChar()
	case '/':
		tok = Token{Type: DIVIDE, Literal: "/", Pos: l.position}
		l.readChar()
	case '(':
		tok = Token{Type: LPAREN, Literal: "(", Pos: l.position}
		l.readChar()
	case ')':
		tok = Token{Type: RPAREN, Literal: ")", Pos: l.position}
		l.readChar()
	case '^':
		tok = Token{Type: CARET, Literal: "^", Pos: l.position}
		l.readChar()
	case '{':
		tok = Token{Type: LBRACE, Literal: "{", Pos: l.position}
		l.readChar()
	case '}':
		tok = Token{Type: RBRACE, Literal: "}", Pos: l.position}
		l.readChar()
	case '[':
		tok = Token{Type: LBRACKET, Literal: "[", Pos: l.position}
		l.readChar()
	case ']':
		tok = Token{Type: RBRACKET, Literal: "]", Pos: l.position}
		l.readChar()
	case '=':
		tok = Token{Type: EQUAL, Literal: "=", Pos: l.position}
		l.readChar()
	case '\\':
		cmdName := l.readCommand()
		if cmdName == "sqrt" {
			tok = Token{Type: COMMAND, Literal: cmdName, Pos: l.position - len(cmdName) - 1}
			return tok
		}
		tok = Token{Type: ILLEGAL, Literal: "\\" + cmdName, Pos: l.position}
		return tok
	case 0:
		tok = Token{Type: EOF, Literal: "", Pos: l.position}
	default:
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
			// Parse the number value
			val, err := strconv.ParseFloat(tok.Literal, 64)
			if err != nil {
				tok.Type = ILLEGAL
			} else {
				tok.Value = val
			}
			return tok
		} else if isLetter(l.ch) {
			// Read a single letter variable
			tok.Literal = string(l.ch)
			tok.Type = VARIABLE
			l.readChar()
			// Check if the next character is also a letter (multi-character identifier)
			// If so, mark as ILLEGAL since we only support single-character variables
			if isLetter(l.ch) {
				tok.Type = ILLEGAL
			}
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Pos: l.position}
			l.readChar()
		}
	}

	return tok
}
