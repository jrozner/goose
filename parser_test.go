package goose

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	/*
		file, err := os.Open("doc/example.txt")
		if err != nil {
			t.Fatal(err)
		}*/

	input := `up add column "users", string, {size: 5, null:true} add column "password", text end down add column "users", boolean end`
	lexer := NewLexer(strings.NewReader(input))

	tree, err := Parse(lexer)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	walk(tree, 0)
}

func walk(node *Node, depth int) {
	if node == nil {
		return
	}

	name, _ := Nodes[node.Type]
	for i := 0; i < depth*2; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%s\n", name)

	for _, child := range node.Children {
		switch t := child.(type) {
		case *Node:
			walk(t, depth+1)
		case *Token:
			for i := 0; i < (depth+1)*2; i++ {
				fmt.Print(" ")
			}
			fmt.Printf("%s\n", t.String())
		default:
			fmt.Println("wtf are you?")
		}
	}
}
