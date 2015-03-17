package goose

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	data := `::{,:,:::}:`
	lexer := NewLexer(strings.NewReader(data))

	for {
		tok := lexer.Next()
		fmt.Println(tok.String())

		if tok.Type == EOF {
			break
		}
	}
}
