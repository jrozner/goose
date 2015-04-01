package goose

import (
	"errors"
	"log"
)

var ErrNoMatch = errors.New("no match")

func Parse(lexer *Lexer) (*Node, error) {
	parser := parser{lexer: lexer, tokens: make([]*Token, 0)}
	return parser.parseRoot()
}

type parser struct {
	lexer  *Lexer
	tokens []*Token
}

func (p *parser) next() *Token {
	if len(p.tokens) > 0 {
		token := p.tokens[0]
		p.tokens = p.tokens[1:]
		return token
	}

	return p.lexer.Next()
}

func (p *parser) backup(node *Node) {
	for i := len(node.Children) - 1; i >= 0; i-- {
		switch asserted := node.Children[i].(type) {
		case *Token:
			p.tokens = append([]*Token{asserted}, p.tokens...)
		case *Node:
			p.backup(asserted)
		default:
			log.Println("wtf is going on here")
		}
	}
}

func (p *parser) parseRoot() (*Node, error) {
	node := &Node{Type: NodeRoot}

	ret, err := p.parseUpStatement()
	if err != nil {
		return node, err
	}

	node.Children = append(node.Children, ret)

	ret, err = p.parseDownStatement()
	if err != nil {
		return node, err
	}

	node.Children = append(node.Children, ret)

	token := p.next()
	if token.Type != TokenEOF {
		return node, ErrNoMatch
	}

	node.Children = append(node.Children, token)

	return node, nil
}

func (p *parser) parseUpStatement() (*Node, error) {
	node := &Node{Type: NodeUpStatement}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenUp {
		p.backup(node)
		return node, ErrNoMatch
	}

	for {
		ret, err := p.parseStatement()
		if err != nil {
			break
		}

		node.Children = append(node.Children, ret)
	}

	token = p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenEnd {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseDownStatement() (*Node, error) {
	node := &Node{Type: NodeDownStatement}

	token := p.next()
	if token.Type != TokenDown {
		return node, ErrNoMatch
	}

	node.Children = append(node.Children, token)

	for {
		ret, err := p.parseStatement()
		if err != nil {
			break
		}

		node.Children = append(node.Children, ret)
	}

	token = p.next()
	if token.Type != TokenEnd {
		return node, ErrNoMatch
	}

	node.Children = append(node.Children, token)

	return node, nil
}

func (p *parser) parseStatement() (*Node, error) {
	node := &Node{Type: NodeStatement}

	child, err := p.parseAddColumn()
	if err == nil {
		node.Children = append(node.Children, child)
		return node, nil
	}

	return nil, ErrNoMatch
}

func (p *parser) parseAddColumn() (*Node, error) {
	node := &Node{Type: NodeAddColumn}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenAdd {
		p.backup(node)
		return nil, ErrNoMatch
	}

	token = p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenColumn {
		p.backup(node)
		return nil, ErrNoMatch
	}

	tableName, err := p.parseTableName()
	if err != nil {
		p.backup(node)
		return nil, ErrNoMatch
	}

	node.Children = append(node.Children, tableName)

	token = p.next()
	node.Children = append(node.Children, token)
	if token.Type != TokenComma {
		p.backup(node)
		return nil, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseTableName() (*Node, error) {
	node := &Node{Type: NodeTableName}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenStringLiteral {
		p.backup(node)
		return nil, ErrNoMatch
	}

	return node, nil
}
