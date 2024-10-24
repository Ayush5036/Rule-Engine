package parser

import (
    "fmt"
    "strconv"
    "strings"
    "unicode"

    "github.com/Ayush/rule-engine/internal/ast"
)

type TokenType int

const (
    OPERATOR TokenType = iota
    FIELD
    VALUE
    LPAREN
    RPAREN
    COMPARISON
    AND
    OR
)

type Token struct {
    Type  TokenType
    Value string
}

type Parser struct {
    tokens  []Token
    current int
}

// Tokenize function updated for case-insensitivity and trimming
func Tokenize(input string) ([]Token, error) {
    var tokens []Token
    input = strings.TrimSpace(input)

    for i := 0; i < len(input); i++ {
        switch {
        case unicode.IsSpace(rune(input[i])):
            continue

        case input[i] == '(':
            tokens = append(tokens, Token{Type: LPAREN, Value: "("})

        case input[i] == ')':
            tokens = append(tokens, Token{Type: RPAREN, Value: ")"})

        case input[i] == '\'':
            j := i + 1
            for j < len(input) && input[j] != '\'' {
                j++
            }
            if j >= len(input) {
                return nil, fmt.Errorf("unterminated string literal")
            }
            tokens = append(tokens, Token{Type: VALUE, Value: strings.TrimSpace(input[i+1:j])})
            i = j

        case unicode.IsLetter(rune(input[i])):
            j := i
            for j < len(input) && (unicode.IsLetter(rune(input[j])) || unicode.IsDigit(rune(input[j]))) {
                j++
            }
            word := strings.TrimSpace(input[i:j])
            switch strings.ToUpper(word) { // Case-insensitive comparison for AND/OR
            case "AND":
                tokens = append(tokens, Token{Type: AND, Value: "AND"})
            case "OR":
                tokens = append(tokens, Token{Type: OR, Value: "OR"})
            default:
                tokens = append(tokens, Token{Type: FIELD, Value: strings.ToLower(word)}) // Convert fields to lowercase
            }
            i = j - 1

        case unicode.IsDigit(rune(input[i])):
            j := i
            for j < len(input) && (unicode.IsDigit(rune(input[j])) || input[j] == '.') {
                j++
            }
            tokens = append(tokens, Token{Type: VALUE, Value: strings.TrimSpace(input[i:j])})
            i = j - 1

        case input[i] == '>' || input[i] == '<' || input[i] == '=':
            tokens = append(tokens, Token{Type: COMPARISON, Value: string(input[i])})

        default:
            return nil, fmt.Errorf("unexpected character: %c", input[i])
        }
    }

    return tokens, nil
}

func (p *Parser) Parse(tokens []Token) (*ast.Node, error) {
    p.tokens = tokens
    p.current = 0
    node, err := p.parseExpression()
    if err != nil {
        return nil, err
    }

    // Ensure there are no leftover tokens after parsing
    if !p.isAtEnd() {
        return nil, fmt.Errorf("unexpected tokens remaining after parsing")
    }

    return node, nil
}

func (p *Parser) parseExpression() (*ast.Node, error) {
    if p.match(LPAREN) {
        expr, err := p.parseParenExpression()
        if err != nil {
            return nil, err
        }

        return expr, nil
    }

    if p.match(FIELD) {
        field := p.previous().Value
        if !p.match(COMPARISON) {
            return nil, fmt.Errorf("expected comparison operator after field")
        }
        operator := p.previous().Value
        if !p.match(VALUE) {
            return nil, fmt.Errorf("expected value after comparison operator")
        }
        value := p.previous().Value

        node := &ast.Node{
            Type:     "operand",
            Field:    field,
            Operator: operator,
            Value:    parseValue(value),
        }

        if p.match(AND) || p.match(OR) {
            op := p.previous().Value
            rightExpr, err := p.parseExpression()
            if err != nil {
                return nil, err
            }

            return &ast.Node{
                Type:     "operator",
                Operator: op,
                Children: []*ast.Node{node, rightExpr},
            }, nil
        }

        return node, nil
    }

    return nil, fmt.Errorf("unexpected token")
}

func (p *Parser) parseParenExpression() (*ast.Node, error) {
    expr, err := p.parseExpression()
    if err != nil {
        return nil, err
    }

    if !p.match(RPAREN) {
        return nil, fmt.Errorf("expected closing parenthesis")
    }

    return expr, nil
}

func (p *Parser) advance() {
    if p.current < len(p.tokens) {
        p.current++
    }
}

func (p *Parser) match(tokenType TokenType) bool {
    if p.isAtEnd() || p.tokens[p.current].Type != tokenType {
        return false
    }
    p.current++
    return true
}

func (p *Parser) previous() Token {
    return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
    return p.current >= len(p.tokens)
}

func parseValue(value string) interface{} {
    if i, err := strconv.Atoi(value); err == nil {
        return i
    }

    if f, err := strconv.ParseFloat(value, 64); err == nil {
        return f
    }

    return value
}

func ParseRule(ruleString string) (*ast.Node, error) {
    tokens, err := Tokenize(ruleString)
    if err != nil {
        return nil, fmt.Errorf("tokenization error: %w", err)
    }

    parser := &Parser{}
    ast, err := parser.Parse(tokens)
    if err != nil {
        return nil, fmt.Errorf("parsing error: %w", err)
    }

    return ast, nil
}
