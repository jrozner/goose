package goose

import (
	"fmt"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	file, err := os.Open("doc/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	lexer := NewLexer(file)

	for {
		tok := lexer.Next()
		fmt.Printf("[%d,%d) %d %s\n", tok.Start, tok.Stop, tok.Type, tok.String())

		if tok.Type == TokenEOF || tok.Type == TokenErr {
			break
		}
	}
}
