package goose

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Lexer struct {
	input    *bufio.Reader
	buffer   []rune
	position int
	tokens   chan *Token
}

// NewLexer returns a pointer to a new Lexer
func NewLexer(input io.Reader) *Lexer {
	lexer := &Lexer{
		input:  bufio.NewReader(input),
		buffer: make([]rune, 0),
		tokens: make(chan *Token, 2),
	}

	go lexer.run()
	return lexer
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

func (l *Lexer) emit(start, stop int, tokenType TokenType, raw []rune, value interface{}) {
	l.tokens <- &Token{
		Start: start,
		Stop:  stop,
		Type:  tokenType,
		Raw:   raw,
		Value: value,
	}
}

func (l *Lexer) skipWhitespace() {
	var (
		raw   = make([]rune, 0)
		start = l.position
	)

	for {
		ch, err := l.peek()
		if err != nil {
			if err == io.EOF {
				l.emit(start, l.position, EOF, []rune{}, err.Error())
			}

			l.emit(start, l.position, Err, raw, err)
			return
		}

		if !unicode.IsSpace(ch) {
			break
		}

		l.next()
		raw = append(raw, ch)
	}
}

func (l *Lexer) consumeString() {
	var (
		raw   = make([]rune, 0)
		start = l.position
	)

	ch, err := l.next()
	if err != nil {
		l.emit(start, l.position, Err, raw, err)
		return
	}

	raw = append(raw, ch)

	for {
		ch, err = l.peek()
		if err != nil {
			l.emit(start, l.position, Err, raw, err)
			return
		}

		switch ch {
		case '"':
			l.next()
			raw = append(raw, ch)
			l.emit(start, l.position, String, raw, string(raw))
			return
		default:
			l.next()
			raw = append(raw, ch)
			continue
		}
	}
}

func (l *Lexer) run() {
	for {
		l.skipWhitespace()

		ch, err := l.peek()
		if err != nil {
			if err == io.EOF {
				l.emit(l.position, l.position, EOF, []rune{}, err.Error())
				break
			}

			l.emit(l.position, l.position, Err, []rune{}, err.Error())
		}

		switch {
		case ch == '{':
			start := l.position
			l.next()
			l.emit(start, l.position, LeftBrace, []rune{ch}, '{')
		case ch == '}':
			start := l.position
			l.next()
			l.emit(start, l.position, RightBrace, []rune{ch}, '}')
		case ch == ',':
			start := l.position
			l.next()
			l.emit(start, l.position, Comma, []rune{ch}, ',')
		case ch == ':':
			start := l.position
			l.next()
			l.emit(start, l.position, Colon, []rune{ch}, ':')
		case ch == '"':
			l.consumeString()
		case unicode.IsLetter(ch):
			l.consumeKeyword()
		case unicode.IsNumber(ch), ch == '-':
			//l.consumeNumber()
		}
	}

	close(l.tokens)
}

func (l *Lexer) consumeKeyword() {
	var (
		start int
		raw   = make([]rune, 0)
	)

	for {
		ch, err := l.peek()
		if err != nil {
			if err == io.EOF {
				l.emit(l.position, l.position, EOF, []rune{}, err.Error())
				break
			}

			l.emit(l.position, l.position, Err, []rune{}, err.Error())
			return
		}

		switch {
		case unicode.IsLower(ch), ch == '_':
			l.next()
			raw = append(raw, ch)
		default:
			if token, ok := keywords[string(raw)]; ok {
				l.emit(start, l.position, token, raw, string(raw))
				return
			}

			l.emit(start, l.position, Err, raw, fmt.Errorf("unexpected %s", string(raw)))
			return
		}
	}
}

func (l *Lexer) Next() *Token {
	return <-l.tokens
}
