package query

import (
	"fmt"
)

// Parser implements a recursive descent parser for query expressions
type Parser struct {
	tokens []Token
	pos    int
}

// Parse converts a query string into an AST
func Parse(query string) (Node, error) {
	tokens, err := Tokenize(query)
	if err != nil {
		return nil, err
	}

	parser := &Parser{tokens: tokens, pos: 0}
	node, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}

	// Ensure we consumed all tokens (except EOF)
	if !parser.isAtEnd() {
		return nil, fmt.Errorf("unexpected token at position %d: %s", parser.current().Pos, parser.current())
	}

	return node, nil
}

// parseExpression handles OR at the top level
// expression := term (OR term)*
func (p *Parser) parseExpression() (Node, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(TokenOR) {
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		node = &BooleanNode{Op: "OR", Left: node, Right: right}
	}

	return node, nil
}

// parseTerm handles AND
// term := factor (AND factor)*
func (p *Parser) parseTerm() (Node, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(TokenAND) {
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		node = &BooleanNode{Op: "AND", Left: node, Right: right}
	}

	return node, nil
}

// parseFactor handles NOT and parentheses
// factor := NOT factor | ( expression ) | comparison
func (p *Parser) parseFactor() (Node, error) {
	// Handle NOT
	if p.match(TokenNOT) {
		node, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return &BooleanNode{Op: "NOT", Left: node}, nil
	}

	// Handle parentheses
	if p.match(TokenLeftParen) {
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.match(TokenRightParen) {
			return nil, fmt.Errorf("expected ) at position %d", p.current().Pos)
		}
		return node, nil
	}

	// Handle comparison
	return p.parseComparison()
}

// parseComparison handles field:value, field>value, etc.
// comparison := FIELD (: | > | < | = | !=) VALUE
func (p *Parser) parseComparison() (Node, error) {
	if !p.check(TokenField) {
		return nil, fmt.Errorf("expected field name at position %d, got %s", p.current().Pos, p.current())
	}

	field := p.advance()

	// Expect an operator
	if !p.check(TokenColon) && !p.check(TokenGT) && !p.check(TokenLT) &&
		!p.check(TokenEQ) && !p.check(TokenNE) {
		return nil, fmt.Errorf("expected operator (:, >, <, =, !=) at position %d, got %s", p.current().Pos, p.current())
	}

	operator := p.advance()

	// Expect a value
	if !p.check(TokenValue) {
		return nil, fmt.Errorf("expected value at position %d, got %s", p.current().Pos, p.current())
	}

	value := p.advance()

	return &ComparisonNode{
		Field:    field.Value,
		Operator: operator.Value,
		Value:    value.Value,
	}, nil
}

// Helper methods

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1] // Return EOF
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.pos++
	}
	return p.tokens[p.pos-1]
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().Type == t
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.current().Type == TokenEOF
}
