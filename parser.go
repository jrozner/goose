package goose

import "fmt"

type Parser func(*Lexer) (*Node, error)

func Parse(lexer *Lexer) (*Node, error) {
	return parseRoot(lexer)
}

func parseRoot(lexer *Lexer) (*Node, error) {
	node := new(Node)

	ret, err := parseUpStatement(lexer)
	if err != nil {
		return node, err
	}

	node.Children = append(node.Children, ret)

	ret, err = parseDownStatement(lexer)
	if err != nil {
		return node, err
	}

	node.Children = append(node.Children, ret)

	return node, nil
}

func parseUpStatement(lexer *Lexer) (*Node, error) {
	node := new(Node)

	token := lexer.Next()
	if token.Type != TokenUp {
		return node, fmt.Errorf("unexpected token: %d", token.Type)
	}

	node.Children = append(node.Children, token)

	ret, err := parseStatement(lexer)
	if err != nil {
		return node, err
	}

	node.Children = append(node.Children, ret)

	return node, nil
}

func parseStatement(lexer *Lexer) (*Node, error) {
	return nil, nil
}

func parseDownStatement(lexer *Lexer) (*Node, error) {
	return nil, nil
}
