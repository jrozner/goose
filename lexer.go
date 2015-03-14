package goose

import (
	"bufio"
	"io"
)

type lexerFunc func(*Lexer) lexerFunc

type Lexer struct {
	input    *bufio.Reader
	buffer   []rune
	position int
	tokens   chan *Token
}

// NewLexer returns a pointer to a new Lexer
func NewLexer(input io.Reader) *Lexer {
	return &Lexer{
		input:  bufio.NewReader(input),
		buffer: make([]rune, 0),
		tokens: make(chan *Token, 2),
	}
}

// peek returns the next rune but does not consume
func (l *Lexer) peek() (rune, error) {
	if len(l.buffer) > 0 {
		return l.buffer[0], nil
	}

	ch, _, err := l.input.ReadRune()
	if err != nil {
		return ch, err
	}

	l.buffer = append(l.buffer, ch)

	return l.buffer[0], nil
}

// next reads the next rune from the input
func (l *Lexer) next() (rune, error) {
	if len(l.buffer) > 0 {
		ch := l.buffer[0]
		l.buffer = l.buffer[1:]
		l.position++
		return ch, nil
	}

	ch, _, err := l.input.ReadRune()
	if err != nil {
		return ch, err
	}

	l.position++
	return ch, nil
}

func (l *Lexer) emit(token *Token) {
	l.tokens <- token
}

func (l *Lexer) skipWhitespace() {
}

func (l *Lexer) Next() *Token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			// TODO: start up the state machine
		}
	}
}
