package goose

import (
	"bufio"
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
			l.emit(start, l.position, Err, raw, err)
			return
		}

		if !unicode.IsSpace(ch) {
			break
		}

		// this should never actually happen because peek should already have
		// read from the stream if needed and pushed the next rune into the
		// buffer which is impossible to result in an error when reading from
		ch, err = l.next()
		if err != nil {
			l.emit(start, l.position, Err, raw, err)
			return
		}

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
			ch, err = l.next()
			if err != nil {
				l.emit(start, l.position, Err, raw, err)
				return
			}

			raw = append(raw, ch)
			l.emit(start, l.position, String, raw, string(raw))
			return
		default:
			ch, err = l.next()
			if err != nil {
				l.emit(start, l.position, Err, raw, err)
				return
			}

			raw = append(raw, ch)
			continue
		}
	}
}

func (l *Lexer) run() {
	for {
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
			ch, err = l.next()
			if err != nil {
				l.emit(l.position, l.position, Err, []rune{}, err)
				return
			}
			l.emit(start, l.position, LeftBrace, []rune{ch}, '{')
		case ch == '}':
			start := l.position
			ch, err = l.next()
			if err != nil {
				l.emit(l.position, l.position, Err, []rune{}, err)
				return
			}
			l.emit(start, l.position, RightBrace, []rune{ch}, '}')
		case ch == ',':
			start := l.position
			ch, err = l.next()
			if err != nil {
				l.emit(l.position, l.position, Err, []rune{}, err)
				return
			}
			l.emit(start, l.position, Comma, []rune{ch}, ',')
		case ch == ':':
			start := l.position
			ch, err = l.next()
			if err != nil {
				l.emit(l.position, l.position, Err, []rune{}, err)
				return
			}
			l.emit(start, l.position, Colon, []rune{ch}, ':')
		case ch == '"':
			l.consumeString()
		case unicode.IsLetter(ch):
			//l.consumeKeyword()
		case unicode.IsNumber(ch), ch == '-':
			//l.consumeNumber()
		}
	}

	close(l.tokens)
}

func (l *Lexer) Next() *Token {
	return <-l.tokens
}
