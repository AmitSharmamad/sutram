package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"algo/src/shared"
)

// TokenType represents different types of tokens in our grammar
type TokenType int

const (
	NUMBER TokenType = iota
	IDENTIFIER
	OPERATOR
	LPAREN
	RPAREN
	COMMA
	EOF
)

// Token represents a token in our lexical analysis
type Token struct {
	Type  TokenType
	Value string
}

// Lexer converts input string into tokens
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips any whitespace
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = Token{Type: LPAREN, Value: "("}
	case ')':
		tok = Token{Type: RPAREN, Value: ")"}
	case ',':
		tok = Token{Type: COMMA, Value: ","}
	case '+', '-', '*', '/', '%', '>', '<', '=', '!':
		// Handle operators
		if l.ch == '>' && l.peekChar() == '=' {
			tok = Token{Type: OPERATOR, Value: ">="}
			l.readChar()
		} else if l.ch == '<' && l.peekChar() == '=' {
			tok = Token{Type: OPERATOR, Value: "<="}
			l.readChar()
		} else if l.ch == '=' && l.peekChar() == '=' {
			tok = Token{Type: OPERATOR, Value: "=="}
			l.readChar()
		} else if l.ch == '!' && l.peekChar() == '=' {
			tok = Token{Type: OPERATOR, Value: "!="}
			l.readChar()
		} else {
			tok = Token{Type: OPERATOR, Value: string(l.ch)}
		}
	case 0:
		tok = Token{Type: EOF, Value: ""}
	default:
		if unicode.IsDigit(rune(l.ch)) {
			return l.readNumber()
		} else if unicode.IsLetter(rune(l.ch)) || l.ch == '_' {
			return l.readIdentifier()
		} else {
			tok = Token{Type: EOF, Value: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// readNumber reads a number token
func (l *Lexer) readNumber() Token {
	position := l.position
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}

	// Handle decimal point
	if l.ch == '.' && unicode.IsDigit(rune(l.peekChar())) {
		l.readChar() // consume the '.'
		for unicode.IsDigit(rune(l.ch)) {
			l.readChar()
		}
	}

	return Token{Type: NUMBER, Value: l.input[position:l.position]}
}

// readIdentifier reads an identifier token
func (l *Lexer) readIdentifier() Token {
	position := l.position
	for unicode.IsLetter(rune(l.ch)) || unicode.IsDigit(rune(l.ch)) || l.ch == '_' {
		l.readChar()
	}
	value := l.input[position:l.position]

	// Check for logical operators
	if strings.ToUpper(value) == "AND" || strings.ToUpper(value) == "OR" || strings.ToUpper(value) == "NOT" {
		return Token{Type: OPERATOR, Value: strings.ToUpper(value)}
	}

	return Token{Type: IDENTIFIER, Value: value}
}

// AST Node types
type NodeType int

const (
	NUMBER_NODE NodeType = iota
	VARIABLE_NODE
	BINARY_OP_NODE
	UNARY_OP_NODE
	FUNCTION_CALL_NODE
)

// Node represents a node in the AST
type Node interface {
	Type() NodeType
	String() string
	Evaluate(ctx *EvaluationContext) (float64, error)
}

// NumberNode represents a numeric literal
type NumberNode struct {
	Value float64
}

func (n *NumberNode) Type() NodeType {
	return NUMBER_NODE
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n *NumberNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	return n.Value, nil
}

// VariableNode represents a variable reference
type VariableNode struct {
	Name string
}

func (n *VariableNode) Type() NodeType {
	return VARIABLE_NODE
}

func (n *VariableNode) String() string {
	return n.Name
}

func (n *VariableNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	if value, exists := ctx.Variables[n.Name]; exists {
		return value, nil
	}
	return 0, fmt.Errorf("variable not found: %s", n.Name)
}

// GetName returns the name of the variable
func (n *VariableNode) GetName() string {
    return n.Name
}

// BinaryOpNode represents a binary operation
type BinaryOpNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *BinaryOpNode) Type() NodeType {
	return BINARY_OP_NODE
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Operator, n.Right.String())
}

func (n *BinaryOpNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	left, err := n.Left.Evaluate(ctx)
	if err != nil {
		return 0, err
	}

	right, err := n.Right.Evaluate(ctx)
	if err != nil {
		return 0, err
	}

	switch n.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, errors.New("division by zero")
		}
		return left / right, nil
	case "%":
		if right == 0 {
			return 0, errors.New("modulo by zero")
		}
		return math.Mod(left, right), nil
	case ">":
		if left > right {
			return 1, nil
		}
		return 0, nil
	case "<":
		if left < right {
			return 1, nil
		}
		return 0, nil
	case ">=":
		if left >= right {
			return 1, nil
		}
		return 0, nil
	case "<=":
		if left <= right {
			return 1, nil
		}
		return 0, nil
	case "==":
		if left == right {
			return 1, nil
		}
		return 0, nil
	case "!=":
		if left != right {
			return 1, nil
		}
		return 0, nil
	case "AND":
		if left != 0 && right != 0 {
			return 1, nil
		}
		return 0, nil
	case "OR":
		if left != 0 || right != 0 {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", n.Operator)
	}
}

// UnaryOpNode represents a unary operation
type UnaryOpNode struct {
	Operator string
	Operand  Node
}

func (n *UnaryOpNode) Type() NodeType {
	return UNARY_OP_NODE
}

func (n *UnaryOpNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Operator, n.Operand.String())
}

func (n *UnaryOpNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	operand, err := n.Operand.Evaluate(ctx)
	if err != nil {
		return 0, err
	}

	switch n.Operator {
	case "-":
		return -operand, nil
	case "+":
		return operand, nil
	case "NOT":
		if operand == 0 {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown unary operator: %s", n.Operator)
	}
}

// FunctionCallNode represents a function call
type FunctionCallNode struct {
	Name      string
	Arguments []Node
}

func (n *FunctionCallNode) Type() NodeType {
	return FUNCTION_CALL_NODE
}

func (n *FunctionCallNode) String() string {
	args := make([]string, len(n.Arguments))
	for i, arg := range n.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", n.Name, strings.Join(args, ", "))
}


// EMA(close, 50, 0)

func (n *FunctionCallNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	args := make([]any, len(n.Arguments))
	for i, arg := range n.Arguments {
        if varNode, isVariable := arg.(*VariableNode); isVariable {
            args[i] = varNode.Name
            continue
        }

		value, err := arg.Evaluate(ctx)
		if err != nil {
			return 0, err
		}

		args[i] = value
	}

	calculator := GetIndFactory(n.Name)
	if calculator == nil {
		return 0, fmt.Errorf("unknown indicator: %s", n.Name)
	}

	result := calculator.Calculate(args, ctx)
	return result, nil
}

// Parser parses tokens into an AST
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

// NewParser creates a new parser
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken() // Initialize currentToken
	p.nextToken() // Initialize peekToken
	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseExpression parses an expression
func (p *Parser) ParseExpression() (Node, error) {
	return p.parseLogicalExpression()
}

// parseLogicalExpression parses logical expressions (AND, OR)
func (p *Parser) parseLogicalExpression() (Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == OPERATOR &&
		(p.currentToken.Value == "AND" || p.currentToken.Value == "OR") {
		operator := p.currentToken.Value
		p.nextToken()

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// parseComparison parses comparison expressions
func (p *Parser) parseComparison() (Node, error) {
	left, err := p.parseArithmeticExpression()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type == OPERATOR &&
		(p.currentToken.Value == ">" || p.currentToken.Value == "<" ||
			p.currentToken.Value == ">=" || p.currentToken.Value == "<=" ||
			p.currentToken.Value == "==" || p.currentToken.Value == "!=") {
		operator := p.currentToken.Value
		p.nextToken()

		right, err := p.parseArithmeticExpression()
		if err != nil {
			return nil, err
		}

		return &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}, nil
	}

	return left, nil
}

// parseArithmeticExpression parses arithmetic expressions (+ and -)
func (p *Parser) parseArithmeticExpression() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == OPERATOR &&
		(p.currentToken.Value == "+" || p.currentToken.Value == "-") {
		operator := p.currentToken.Value
		p.nextToken()

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// parseTerm parses terms (* and /)
func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == OPERATOR &&
		(p.currentToken.Value == "*" || p.currentToken.Value == "/" || p.currentToken.Value == "%") {
		operator := p.currentToken.Value
		p.nextToken()

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// parseFactor parses factors
func (p *Parser) parseFactor() (Node, error) {
	if p.currentToken.Type == OPERATOR &&
		(p.currentToken.Value == "+" || p.currentToken.Value == "-" || p.currentToken.Value == "NOT") {
		operator := p.currentToken.Value
		p.nextToken()

		operand, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return &UnaryOpNode{
			Operator: operator,
			Operand:  operand,
		}, nil
	}

	return p.parsePrimary()
}

// parsePrimary parses primary expressions
func (p *Parser) parsePrimary() (Node, error) {
	switch p.currentToken.Type {
	case NUMBER:
		value, err := strconv.ParseFloat(p.currentToken.Value, 64)
		if err != nil {
			return nil, err
		}
		node := &NumberNode{Value: value}
		p.nextToken()
		return node, nil

	case IDENTIFIER:
		identifier := p.currentToken.Value
		p.nextToken()

		// Check if it's a function call
		if p.currentToken.Type == LPAREN {
			p.nextToken() // consume '('
			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}
			return &FunctionCallNode{
				Name:      identifier,
				Arguments: args,
			}, nil
		}

		// It's a variable reference
		return &VariableNode{Name: strings.ToUpper(identifier)}, nil

	case LPAREN:
		p.nextToken() // consume '('
		expr, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}

		if p.currentToken.Type != RPAREN {
			return nil, errors.New("expected ')'")
		}
		p.nextToken() // consume ')'

		return expr, nil

	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Value)
	}
}

// parseArguments parses function arguments
func (p *Parser) parseArguments() ([]Node, error) {
	var args []Node

	// Handle empty argument list
	if p.currentToken.Type == RPAREN {
		p.nextToken() // consume ')'
		return args, nil
	}

	// Parse first argument
	arg, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	args = append(args, arg)

	// Parse remaining arguments
	for p.currentToken.Type == COMMA {
		p.nextToken() // consume ','

		arg, err := p.parseLogicalExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	if p.currentToken.Type != RPAREN {
		return nil, errors.New("expected ')' after function arguments")
	}
	p.nextToken() // consume ')'

	return args, nil
}

// TradingFunctions defines the interface for trading-specific functions
type TradingFunctions interface {
	SMA(value float64, period int) (float64, error)
	EMA(value float64, period int) (float64, error)
	RSI(period int) (float64, error)
	MACD(fastPeriod, slowPeriod, signalPeriod int) (float64, error)
	BollingerBandUpper(period int, stdDev float64) (float64, error)
	BollingerBandLower(period int, stdDev float64) (float64, error)
	BollingerBandMiddle(period int, stdDev float64) (float64, error)
	ATR(period int) (float64, error)
	Stochastic(kPeriod, kSlowing, dPeriod int) (float64, error)
	CrossAbove(value1, value2 float64) (float64, error)
	CrossBelow(value1, value2 float64) (float64, error)
	PercentChange(current float64, periods int) (float64, error)
	Highest(value float64, period int) (float64, error)
	Lowest(value float64, period int) (float64, error)
}

// EvaluationContext holds the context for expression evaluation
type EvaluationContext struct {
	Data *shared.StockData

	Variables map[string]float64
	// Functions TradingFunctions
}

// NewEvaluationContext creates a new evaluation context
func NewEvaluationContext(data *shared.StockData) *EvaluationContext {
	return &EvaluationContext{
		Data: data,
		// Variables: vars,
		// Functions: funcs,
	}
}

// ExpressionBuilder allows building expressions programmatically
type ExpressionBuilder struct {
	root Node
}

// NewExpressionBuilder creates a new expression builder
func NewExpressionBuilder() *ExpressionBuilder {
	return &ExpressionBuilder{}
}

// Number creates a number node
func (b *ExpressionBuilder) Number(value float64) *NumberNode {
	return &NumberNode{Value: value}
}

// Variable creates a variable node
func (b *ExpressionBuilder) Variable(name string) *VariableNode {
	return &VariableNode{Name: strings.ToUpper(name)}
}

// BinaryOp creates a binary operation node
func (b *ExpressionBuilder) BinaryOp(left Node, operator string, right Node) *BinaryOpNode {
	return &BinaryOpNode{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// UnaryOp creates a unary operation node
func (b *ExpressionBuilder) UnaryOp(operator string, operand Node) *UnaryOpNode {
	return &UnaryOpNode{
		Operator: operator,
		Operand:  operand,
	}
}

// FunctionCall creates a function call node
func (b *ExpressionBuilder) FunctionCall(name string, args ...Node) *FunctionCallNode {
	return &FunctionCallNode{
		Name:      name,
		Arguments: args,
	}
}

// Parse parses an expression string
func Parse(expr string) (Node, error) {
	lexer := NewLexer(expr)
	parser := NewParser(lexer)
	return parser.ParseExpression()
}

// Execute evaluates an expression
func Execute(expr string, ctx *EvaluationContext) bool {
	node, err := Parse(expr)
	if err != nil {
		println("Error in parsing the expression: ", err)
		panic("parse expression error: ALERT ALERT ALERT!!!")
	}

	result, err := node.Evaluate(ctx)
	if err != nil {
		println("Error in evaluating the expression: ", err)
		panic("parse evaluation error: ALERT ALERT ALERT!!!")
	}

	return result > 0
}
