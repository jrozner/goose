package goose

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
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
				l.emit(start, l.position, TokenEOF, []rune{}, err.Error())
			}

			l.emit(start, l.position, TokenErr, raw, err)
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
		l.emit(start, l.position, TokenErr, raw, err)
		return
	}

	raw = append(raw, ch)

	for {
		ch, err = l.peek()
		if err != nil {
			l.emit(start, l.position, TokenErr, raw, err)
			return
		}

		switch ch {
		case '"':
			l.next()
			raw = append(raw, ch)
			l.emit(start, l.position, TokenStringLiteral, raw, string(raw))
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
				l.emit(l.position, l.position, TokenEOF, []rune{}, err.Error())
				break
			}

			l.emit(l.position, l.position, TokenErr, []rune{}, err.Error())
		}

		switch {
		case ch == '{':
			start := l.position
			l.next()
			l.emit(start, l.position, TokenLeftBrace, []rune{ch}, '{')
		case ch == '}':
			start := l.position
			l.next()
			l.emit(start, l.position, TokenRightBrace, []rune{ch}, '}')
		case ch == ',':
			start := l.position
			l.next()
			l.emit(start, l.position, TokenComma, []rune{ch}, ',')
		case ch == ':':
			start := l.position
			l.next()
			l.emit(start, l.position, TokenColon, []rune{ch}, ':')
		case ch == '"':
			l.consumeString()
		case unicode.IsLetter(ch):
			l.consumeKeyword()
		case unicode.IsNumber(ch), ch == '-':
			l.consumeNumber()
		default:
			l.emit(l.position, l.position, TokenErr, []rune{ch}, errors.New("unexpected input"))
		}
	}

	close(l.tokens)
}

func (l *Lexer) consumeNumber() {
	var (
		start      = l.position
		raw        = make([]rune, 0)
		isFloat    bool
		isNegative bool
	)

	for {
		ch, err := l.peek()
		if err != nil {
			if err == io.EOF {
				l.emit(start, l.position, TokenEOF, []rune{}, err.Error())
				break
			}

			l.emit(start, l.position, TokenErr, []rune{}, err.Error())
			return
		}

		switch {
		case ch == '-':
			if isNegative == true {
				l.emit(start, l.position, TokenErr, raw, errors.New("unexpected input"))
				return
			}

			isNegative = true
			l.next()
			raw = append(raw, ch)
		case unicode.IsNumber(ch):
			l.next()
			raw = append(raw, ch)
		case ch == 'e', ch == 'E', ch == '.':
			if isFloat == true {
				l.emit(start, l.position, TokenErr, raw, errors.New("unexpected input"))
				return
			}

			isFloat = true
			l.next()
			raw = append(raw, ch)
		default:
			if isFloat {
				num, err := strconv.ParseFloat(string(raw), 64)
				if err != nil {
					l.emit(start, l.position, TokenErr, raw, err)
					return
				}

				l.emit(start, l.position, TokenFloatLiteral, raw, num)
			} else {
				num, err := strconv.ParseInt(string(raw), 10, 64)
				if err != nil {
					l.emit(start, l.position, TokenErr, raw, err)
					return
				}

				l.emit(start, l.position, TokenIntegerLiteral, raw, num)
			}

			return
		}
	}
}

func (l *Lexer) consumeKeyword() {
	var (
		start = l.position
		raw   = make([]rune, 0)
	)

	for {
		ch, err := l.peek()
		if err != nil && err != io.EOF {
			l.emit(l.position, l.position, TokenErr, []rune{}, err.Error())
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

			l.emit(start, l.position, TokenErr, raw, fmt.Errorf("unexpected %s", string(raw)))
			return
		}
	}
}

func (l *Lexer) Next() *Token {
	return <-l.tokens
}
