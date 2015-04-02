package goose

import "errors"

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
			panic("unknown type encountered")
		}
	}
}

func (p *parser) backupToken(token *Token) {
	p.tokens = append([]*Token{token}, p.tokens...)
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

	return node, ErrNoMatch
}

func (p *parser) parseAddColumn() (*Node, error) {
	node := &Node{Type: NodeAddColumn}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenAdd {
		p.backup(node)
		return node, ErrNoMatch
	}

	token = p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenColumn {
		p.backup(node)
		return node, ErrNoMatch
	}

	tableName, err := p.parseTableName()
	if err != nil {
		p.backup(node)
		return node, ErrNoMatch
	}

	node.Children = append(node.Children, tableName)

	token = p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenComma {
		p.backup(node)
		return node, ErrNoMatch
	}

	dataType, err := p.parseDataType()
	node.Children = append(node.Children, dataType)

	if err != nil {
		p.backup(node)
		return node, ErrNoMatch
	}

	// Not matching a comma here isn't a parse failure. It just means there's no
	// options block specified
	token = p.next()
	if token.Type != TokenComma {
		p.backupToken(token)
		return node, nil
	}

	node.Children = append(node.Children, token)

	optionsBlock, err := p.parseOptionsBlock()
	node.Children = append(node.Children, optionsBlock)

	if err != nil {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseTableName() (*Node, error) {
	node := &Node{Type: NodeTableName}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenStringLiteral {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseOptionsBlock() (*Node, error) {
	node := &Node{Type: NodeOptionsBlock}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenLeftBrace {
		p.backup(node)
		return node, ErrNoMatch
	}

	option, err := p.parseOption()
	node.Children = append(node.Children, option)

	if err != nil {
		return node, ErrNoMatch
	}

	for {
		token = p.next()
		if token.Type != TokenComma {
			p.backupToken(token)
			break
		}

		node.Children = append(node.Children, token)

		option, err = p.parseOption()
		node.Children = append(node.Children, option)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}
	}

	token = p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenRightBrace {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseOption() (*Node, error) {
	node := &Node{Type: NodeOption}

	// default
	token := p.next()
	if token.Type == TokenDefault {
		node.Children = append(node.Children, token)

		token = p.next()
		node.Children = append(node.Children, token)

		if token.Type != TokenColon {
			p.backup(node)
			return node, ErrNoMatch
		}

		defaultValue, err := p.parseDefaultValue()
		node.Children = append(node.Children, defaultValue)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}

		if defaultValue.Type != NodeDefaultValue {
			p.backup(node)
			return node, ErrNoMatch
		}

		return node, nil
	} else {
		p.backupToken(token)
	}

	// null
	token = p.next()
	if token.Type == TokenNull {
		node.Children = append(node.Children, token)

		token = p.next()
		node.Children = append(node.Children, token)

		if token.Type != TokenColon {
			p.backup(node)
			return node, ErrNoMatch
		}

		boolean, err := p.parseBoolean()
		node.Children = append(node.Children, boolean)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}

		if boolean.Type != NodeBoolean {
			p.backup(node)
			return node, ErrNoMatch
		}

		return node, nil
	} else {
		p.backupToken(token)
	}

	// size
	token = p.next()
	if token.Type == TokenSize {
		node.Children = append(node.Children, token)

		token = p.next()
		node.Children = append(node.Children, token)

		if token.Type != TokenColon {
			p.backup(node)
			return node, ErrNoMatch
		}

		integer, err := p.parseInteger()
		node.Children = append(node.Children, integer)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}

		if integer.Type != NodeInteger {
			p.backup(node)
			return node, ErrNoMatch
		}

		return node, nil
	} else {
		p.backupToken(token)
	}

	// precision
	token = p.next()
	if token.Type == TokenPrecision {
		node.Children = append(node.Children, token)

		token = p.next()
		node.Children = append(node.Children, token)

		if token.Type != TokenColon {
			p.backup(node)
			return node, ErrNoMatch
		}

		integer, err := p.parseInteger()
		node.Children = append(node.Children, integer)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}

		if integer.Type != NodeInteger {
			p.backup(node)
			return node, ErrNoMatch
		}

		return node, nil
	} else {
		p.backupToken(token)
	}

	// scale
	token = p.next()
	if token.Type == TokenScale {
		node.Children = append(node.Children, token)

		token = p.next()
		node.Children = append(node.Children, token)

		if token.Type != TokenColon {
			p.backup(node)
			return node, ErrNoMatch
		}

		integer, err := p.parseInteger()
		node.Children = append(node.Children, integer)

		if err != nil {
			p.backup(node)
			return node, ErrNoMatch
		}

		if integer.Type != NodeInteger {
			p.backup(node)
			return node, ErrNoMatch
		}

		return node, nil
	} else {
		p.backupToken(token)
	}

	// We made it to the end and none of the options matched
	return node, ErrNoMatch
}

func (p *parser) parseDataType() (*Node, error) {
	node := &Node{Type: NodeDataType}

	token := p.next()

	switch token.Type {
	case TokenBinary, TokenBoolean, TokenDate, TokenDatetime, TokenDecimal, TokenFloat, TokenInteger, TokenPrimaryKey, TokenReferences, TokenString, TokenText, TokenTime, TokenTimestamp:
		node.Children = append(node.Children, token)
		return node, nil
	default:
		p.backupToken(token)
		return node, ErrNoMatch
	}

	panic("unknown type encountered")
}

func (p *parser) parseDefaultValue() (*Node, error) {
	node := &Node{Type: NodeDefaultValue}

	token := p.next()

	switch token.Type {
	case TokenTrue, TokenFalse, TokenFloatLiteral, TokenIntegerLiteral, TokenNull, TokenString:
		node.Children = append(node.Children, token)
		return node, nil
	default:
		p.backupToken(token)
		return node, ErrNoMatch
	}

	panic("unknown type encountered")
}

func (p *parser) parseBoolean() (*Node, error) {
	node := &Node{Type: NodeBoolean}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenTrue && token.Type != TokenFalse {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseFloat() (*Node, error) {
	node := &Node{Type: NodeFloat}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenFloatLiteral {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}

func (p *parser) parseInteger() (*Node, error) {
	node := &Node{Type: NodeInteger}

	token := p.next()
	node.Children = append(node.Children, token)

	if token.Type != TokenIntegerLiteral {
		p.backup(node)
		return node, ErrNoMatch
	}

	return node, nil
}
