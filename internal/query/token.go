package query

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of token
type TokenType int

const (
	TokenField TokenType = iota
	TokenColon
	TokenGT
	TokenLT
	TokenEQ
	TokenNE
	TokenValue
	TokenAND
	TokenOR
	TokenNOT
	TokenLeftParen
	TokenRightParen
	TokenEOF
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

func (t Token) String() string {
	switch t.Type {
	case TokenField:
		return fmt.Sprintf("FIELD(%s)", t.Value)
	case TokenColon:
		return ":"
	case TokenGT:
		return ">"
	case TokenLT:
		return "<"
	case TokenEQ:
		return "="
	case TokenNE:
		return "!="
	case TokenValue:
		return fmt.Sprintf("VALUE(%s)", t.Value)
	case TokenAND:
		return "AND"
	case TokenOR:
		return "OR"
	case TokenNOT:
		return "NOT"
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

// Tokenize converts a query string into tokens
func Tokenize(query string) ([]Token, error) {
	var tokens []Token
	pos := 0
	query = strings.TrimSpace(query)

	for pos < len(query) {
		// Skip whitespace
		if unicode.IsSpace(rune(query[pos])) {
			pos++
			continue
		}

		// Check for operators
		if pos < len(query) {
			switch query[pos] {
			case '(':
				tokens = append(tokens, Token{Type: TokenLeftParen, Value: "(", Pos: pos})
				pos++
				continue
			case ')':
				tokens = append(tokens, Token{Type: TokenRightParen, Value: ")", Pos: pos})
				pos++
				continue
			case ':':
				tokens = append(tokens, Token{Type: TokenColon, Value: ":", Pos: pos})
				pos++
				continue
			case '>':
				tokens = append(tokens, Token{Type: TokenGT, Value: ">", Pos: pos})
				pos++
				continue
			case '<':
				tokens = append(tokens, Token{Type: TokenLT, Value: "<", Pos: pos})
				pos++
				continue
			case '=':
				tokens = append(tokens, Token{Type: TokenEQ, Value: "=", Pos: pos})
				pos++
				continue
			case '!':
				if pos+1 < len(query) && query[pos+1] == '=' {
					tokens = append(tokens, Token{Type: TokenNE, Value: "!=", Pos: pos})
					pos += 2
					continue
				}
			}
		}

		// Read word (field, value, or keyword)
		start := pos
		for pos < len(query) && !unicode.IsSpace(rune(query[pos])) &&
			query[pos] != '(' && query[pos] != ')' &&
			query[pos] != ':' && query[pos] != '>' &&
			query[pos] != '<' && query[pos] != '=' && query[pos] != '!' {
			pos++
		}

		if pos == start {
			return nil, fmt.Errorf("unexpected character at position %d: %c", pos, query[pos])
		}

		word := query[start:pos]
		wordUpper := strings.ToUpper(word)

		// Check for keywords
		switch wordUpper {
		case "AND":
			tokens = append(tokens, Token{Type: TokenAND, Value: word, Pos: start})
		case "OR":
			tokens = append(tokens, Token{Type: TokenOR, Value: word, Pos: start})
		case "NOT":
			tokens = append(tokens, Token{Type: TokenNOT, Value: word, Pos: start})
		default:
			// Determine if this is a field or value based on context
			// If the last token was an operator (:, >, <, =, !=), it's a value
			// Otherwise, it's a field
			if len(tokens) > 0 {
				lastType := tokens[len(tokens)-1].Type
				if lastType == TokenColon || lastType == TokenGT ||
					lastType == TokenLT || lastType == TokenEQ || lastType == TokenNE {
					tokens = append(tokens, Token{Type: TokenValue, Value: word, Pos: start})
				} else {
					tokens = append(tokens, Token{Type: TokenField, Value: word, Pos: start})
				}
			} else {
				tokens = append(tokens, Token{Type: TokenField, Value: word, Pos: start})
			}
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Pos: pos})
	return tokens, nil
}
