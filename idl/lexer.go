package idl

import (
	"fmt"
	"os"
)

type Lexer struct {
	file_name string
	source    string
	line_no   int
	start     int
	current   int
}

func makeLexer(filename string) (Lexer, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return (Lexer{}), err
	}

	return Lexer{
		file_name: filename,
		source:    string(data),
		line_no:   0,
		start:     0,
		current:   0,
	}, nil
}

func (lexer *Lexer) peek(offset int) rune {
	var postion int = lexer.current + offset
	var return_val rune = (0)

	if postion < len(lexer.source) && postion >= 0 {
		return_val = rune(lexer.source[postion])
	}

	return return_val
}

func (lexer *Lexer) Error(cause string) error {
	var line_start int = 0
	for lexer.peek(line_start) != '\n' {
		line_start--
	}
	line_start += lexer.current

	var line_end int = 0
	for lexer.peek(line_end) != '\n' {
		line_end++
	}
	line_end += lexer.current

	return LexerError{
		cause: cause,
		snippet: SourceSnippet{
			file_name: lexer.file_name,
			line:      string(lexer.source[line_start:line_end]),
			start:     line_start + (lexer.start - line_start),
			end:       line_start + (lexer.current - line_start),
		},
	}
}

func (lexer *Lexer) advance() rune {
	lexer.current++
	return rune(lexer.source[lexer.current])
}

func (lexer *Lexer) match(char rune) bool {
	if lexer.current >= len(lexer.source) {
		return false
	}
	if rune(lexer.source[lexer.current]) != char {
		return false
	}

	lexer.current++

	return true
}

func (lexer *Lexer) lexeme() string {
	return lexer.source[lexer.start:lexer.current]
}

func (lexer *Lexer) lexeme_start() int {
	var return_value = lexer.current
	for i := lexer.start; i >= 0; i++ {
		if lexer.source[i] == '\n' {
			return_value = i
			break
		}
	}
	return return_value
}

func (lexer *Lexer) lexeme_end() int {
	return lexer.lexeme_start() + (lexer.current - lexer.start)
}

func (lexer *Lexer) Token(token_type string) Token {
	return makeToken(token_type, lexer.lexeme(), lexer.line_no, lexer.lexeme_start(), lexer.lexeme_end())
}

func (lexer *Lexer) isDigit(current rune) bool {
	return current >= '0' && current <= '9'
}

func (lexer *Lexer) isAlpha(current rune) bool {
	return (current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z') || current == '_'
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.source)
}

func (lexer *Lexer) Char() (Token, error) {
	for lexer.peek(0) != '\'' && !lexer.isAtEnd() {
		if lexer.peek(0) == '\n' {
			lexer.line_no++
		}
		lexer.advance()
	}

	if lexer.isAtEnd() {
		return (Token{}), lexer.Error("Unterminated string")
	}

	lexer.advance()

	return lexer.Token(CharLitteral), nil
}

func (lexer *Lexer) String() (Token, error) {
	for lexer.peek(0) != '"' && !lexer.isAtEnd() {
		if lexer.peek(0) == '\n' {
			lexer.line_no++
		}
		lexer.advance()
	}

	if lexer.isAtEnd() {
		return (Token{}), lexer.Error("Unterminated string")
	}

	lexer.advance()

	return lexer.Token(StringLitteral), nil
}

func (lexer *Lexer) Octal() Token {
	return lexer.Token(Octect)
}

func (lexer *Lexer) Int() Token {
	return lexer.Token(Integer)
}

func (Lexer *Lexer) Identifier() Token {
	return Lexer.Token(Identifier)
}

func (lexer *Lexer) Next() (Token, error) {

	lexer.start = lexer.current

	var current rune = lexer.advance()

	switch current {
	case ' ':
		return Token{}, nil
	case '\t':
		return Token{}, nil
	case '\n':
		lexer.line_no++
		return Token{}, nil
	case ':':
		if lexer.match(':') {
			return lexer.Token(ColonColon), nil
		}
		return lexer.Token(Colon), nil
	case ';':
		return lexer.Token(SemiColon), nil
	case '{':
		return lexer.Token(RightBrace), nil
	case '}':
		return lexer.Token(LeftBrace), nil
	case ',':
		return lexer.Token(Comma), nil
	case '-':
		return lexer.Token(Minus), nil
	case '+':
		return lexer.Token(Plus), nil
	case '=':
		return lexer.Token(Equals), nil
	case '(':
		return lexer.Token(LeftParenthsis), nil
	case ')':
		return lexer.Token(RightParenthsis), nil
	case '<':
		return lexer.Token(LeftAngleBrace), nil
	case '>':
		return lexer.Token(RightAngleBrace), nil
	case '\'':
		return lexer.Char()
	case '"':
		return lexer.String()
	case '\\':
		return lexer.Token(BackSlash), nil
	case '/':
		return lexer.Token(Slash), nil
	case '|':
		if lexer.match('|') {
			return lexer.Token(PipePipe), nil
		}
		return lexer.Token(Pipe), nil
	case '^':
		return lexer.Token(Carot), nil
	case '&':
		if lexer.match('&') {
			return lexer.Token(AmpersandAmpersand), nil
		}
		return lexer.Token(Ampersand), nil
	case '*':
		return lexer.Token(Star), nil
	case '%':
		return lexer.Token(Percent), nil
	case '~':
		return lexer.Token(Tilde), nil
	case '@':
		return lexer.Token(AtSign), nil
	case '#':
		if lexer.match('#') {
			return lexer.Token(HashHash), nil
		}
		return lexer.Token(Hash), nil
	case '!':
		return lexer.Token(Bang), nil
	case '0':
		return lexer.Octal(), nil
	default:
		if lexer.isDigit(current) {
			return lexer.Int(), nil
		} else if lexer.isAlpha(current) {
			return lexer.Identifier(), nil
		}
	}

	var message string = fmt.Sprintf("Unregonised token \"%s\"", lexer.source[lexer.start:lexer.current])
	return Token{}, lexer.Error(message)
}
