package parser

import (
	instrument "algo/internal/modules/instrument/domain/model"
	strategyDto "algo/internal/modules/strategy/application/dto"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
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

// IndicatorInfo represents information about an indicator function call
type IndicatorInfo struct {
	Name      string `json:"name"`
	Arguments []any  `json:"arguments"`
}

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
	if strings.ToUpper(value) == "AND" || strings.ToUpper(value) == "OR" || strings.ToUpper(value) == "NOT" ||
		strings.ToUpper(value) == "CROSSESABOVE" || strings.ToUpper(value) == "CROSSESBELOW" {
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

	var lindicators []IndicatorInfo
	extractFunctionCalls(n.Left, &lindicators)

	calculator := GetIndFactory(lindicators[0].Name)
	uniqueString, err := calculator.GenerateIndicatorKey(lindicators[0].Arguments, nil)
	left := ctx.IndicatorArrayMap[uniqueString][ctx.Index].Value


	if err != nil {
		return 0, err
	}

	var rindicators []IndicatorInfo
	extractFunctionCalls(n.Right, &rindicators)

	rcalculator := GetIndFactory(rindicators[0].Name)
	runiqueString, err := rcalculator.GenerateIndicatorKey(rindicators[0].Arguments, nil)
	right := ctx.IndicatorArrayMap[runiqueString][ctx.Index].Value


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
	case "CROSSESABOVE":
		// Handle existing CROSSESABOVE logic with legacy data structure
		if ctx.Data != nil {
			d := (*ctx.Data)[0 : len(*ctx.Data)-1]
			prevLeft, err := n.Left.Evaluate(NewEvaluationContext(&d))
			if err != nil {
				return 0, nil
			}
			prevRight, err := n.Right.Evaluate(NewEvaluationContext(&d))
			if err != nil {
				return 0, nil
			}
			if (left > right) && (prevLeft <= prevRight) {
				return 1, nil
			}
		}
		// Handle CROSSESABOVE with indicator array map
		if ctx.IndicatorArrayMap != nil && ctx.CurrentTimestamp != nil {
			return n.evaluateCrossesAbove(left, right, ctx)
		}
		return 0, nil
	case "<":
		if left < right {
			return 1, nil
		}
		return 0, nil
	case "CROSSESBELOW":
		// Handle existing CROSSESBELOW logic with legacy data structure
		if ctx.Data != nil {
			d := (*ctx.Data)[0 : len(*ctx.Data)-1]
			prevLeft, err := n.Left.Evaluate(NewEvaluationContext(&d))
			if err != nil {
				return 0, nil
			}
			prevRight, err := n.Right.Evaluate(NewEvaluationContext(&d))
			if err != nil {
				return 0, nil
			}
			if (left < right) && (prevLeft >= prevRight) {
				return 1, nil
			}
		}
		// Handle CROSSESBELOW with indicator array map
		if ctx.IndicatorArrayMap != nil && ctx.CurrentTimestamp != nil {
			return n.evaluateCrossesBelow(left, right, ctx)
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

// evaluateCrossesAbove checks if left crosses above right using indicator array map
func (n *BinaryOpNode) evaluateCrossesAbove(left, right float64, ctx *EvaluationContext) (float64, error) {
	prevLeft, err := n.getPreviousValue(n.Left, ctx)
	if err != nil {
		return 0, nil // Handle gracefully
	}

	prevRight, err := n.getPreviousValue(n.Right, ctx)
	if err != nil {
		return 0, nil // Handle gracefully
	}

	if (left > right) && (prevLeft <= prevRight) {
		return 1, nil
	}
	return 0, nil
}

// evaluateCrossesBelow checks if left crosses below right using indicator array map
func (n *BinaryOpNode) evaluateCrossesBelow(left, right float64, ctx *EvaluationContext) (float64, error) {
	prevLeft, err := n.getPreviousValue(n.Left, ctx)
	if err != nil {
		return 0, nil // Handle gracefully
	}

	prevRight, err := n.getPreviousValue(n.Right, ctx)
	if err != nil {
		return 0, nil // Handle gracefully
	}

	if (left < right) && (prevLeft >= prevRight) {
		return 1, nil
	}
	return 0, nil
}

// getPreviousValue gets the previous value for a node from indicator array map
func (n *BinaryOpNode) getPreviousValue(node Node, ctx *EvaluationContext) (float64, error) {
	if ctx.IndicatorArrayMap == nil || ctx.CurrentTimestamp == nil {
		return 0, fmt.Errorf("indicator array map or timestamp not available")
	}

	// Create a context for previous timestamp evaluation
	prevCtx := &EvaluationContext{
		IndicatorArrayMap: ctx.IndicatorArrayMap,
		CurrentTimestamp:  ctx.PreviousTimestamp,
		PreviousTimestamp: nil, // We don't need to go further back
		Variables:         ctx.Variables,
		Strategy:          ctx.Strategy,
		Results:           ctx.Results,
	}

	return node.Evaluate(prevCtx)
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

func (n *FunctionCallNode) This() FunctionCallNode {
	return *n
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

func (n *FunctionCallNode) Evaluate(ctx *EvaluationContext) (float64, error) {
	// If we have indicator array map, use the new evaluation method
	if ctx.IndicatorArrayMap != nil && ctx.CurrentTimestamp != nil {
		return n.evaluateFromIndicatorMap(ctx)
	}

	// Otherwise, use the existing evaluation method for backward compatibility
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

// evaluateFromIndicatorMap evaluates function calls using the indicator array map
func (n *FunctionCallNode) evaluateFromIndicatorMap(ctx *EvaluationContext) (float64, error) {
	// key, err := n.generateIndicatorKey()
	key := "srrin"

	values, exists := ctx.IndicatorArrayMap[key]
	if !exists {
		return 0, fmt.Errorf("indicator %s not found in array map", key)
	}

	// Find the value at the current timestamp
	for _, val := range values {
		if val.Timestamp.Equal(*ctx.CurrentTimestamp) {
			return val.Value, nil
		}
	}

	return 0, fmt.Errorf("no data found for indicator %s at timestamp %v", key, *ctx.CurrentTimestamp)
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
			p.currentToken.Value == "==" || p.currentToken.Value == "!=" ||
			p.currentToken.Value == "CROSSESABOVE" || p.currentToken.Value == "CROSSESBELOW") {
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

// EvaluationContext holds the context for expression evaluation
type EvaluationContext struct {
	// Legacy support
	Data      *instrument.StockData
	Variables map[string]float64
	Strategy  *strategyDto.StrategyResponse
	Results   map[int64]map[string]interface{}

	// New indicator array map support
	IndicatorArrayMap map[string][]IndicatorValues
	CurrentTimestamp  *time.Time
	PreviousTimestamp *time.Time
	Index             int
}

// NewEvaluationContext creates a new evaluation context (legacy)
func NewEvaluationContext(data *instrument.StockData) *EvaluationContext {
	return &EvaluationContext{
		Data: data,
	}
}

// NewIndicatorMapEvaluationContext creates a new evaluation context for indicator array map
func NewIndicatorMapEvaluationContext(
	indicatorArrayMap map[string][]IndicatorValues,
	currentTimestamp time.Time,
	previousTimestamp *time.Time,
	index int,
) *EvaluationContext {
	return &EvaluationContext{
		IndicatorArrayMap: indicatorArrayMap,
		CurrentTimestamp:  &currentTimestamp,
		PreviousTimestamp: previousTimestamp,
		Index:             index,
		Variables:         make(map[string]float64),
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

// Execute evaluates an expression (legacy)
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

// EvaluateIndicatorExpression evaluates an expression against indicator array map
func EvaluateIndicatorExpression(
	expr string,
	indicatorArrayMap map[string][]IndicatorValues,
	currentTimestamp time.Time,
	previousTimestamp *time.Time,
	index int,
) (bool, error) {
	node, err := Parse(expr)
	if err != nil {
		return false, fmt.Errorf("failed to parse expression: %w", err)
	}

	ctx := NewIndicatorMapEvaluationContext(indicatorArrayMap, currentTimestamp, previousTimestamp, index)

	result, err := node.Evaluate(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	return result > 0, nil
}

// EvaluateConditionAtTimestamp evaluates the expression at a specific timestamp
func EvaluateConditionAtTimestamp(
	expr string,
	indicatorArrayMap map[string][]IndicatorValues,
	timestamp time.Time,
) (bool, error) {
	// Find previous timestamp for cross operations
	var previousTimestamp *time.Time

	// Get all unique timestamps and sort them to find the previous one
	timestampSet := make(map[time.Time]bool)
	for _, values := range indicatorArrayMap {
		for _, val := range values {
			timestampSet[val.Timestamp] = true
		}
	}

	// Convert to slice and sort
	timestamps := make([]time.Time, 0, len(timestampSet))
	for ts := range timestampSet {
		timestamps = append(timestamps, ts)
	}

	// Sort timestamps
	for i := 0; i < len(timestamps)-1; i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i].After(timestamps[j]) {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Find previous timestamp
	for i, ts := range timestamps {
		if ts.Equal(timestamp) && i > 0 {
			previousTimestamp = &timestamps[i-1]
			break
		}
	}

	return EvaluateIndicatorExpression(expr, indicatorArrayMap, timestamp, previousTimestamp, 0)
}

// EvaluateConditionForAllTimestamps evaluates the expression for all timestamps
func EvaluateConditionForAllTimestamps(
	expr string,
	indicatorArrayMap map[string][]IndicatorValues,
) (map[time.Time]bool, error) {
	// Get all unique timestamps
	timestampSet := make(map[time.Time]bool)
	for _, values := range indicatorArrayMap {
		for _, val := range values {
			timestampSet[val.Timestamp] = true
		}
	}

	results := make(map[time.Time]bool)

	for timestamp := range timestampSet {
		result, err := EvaluateConditionAtTimestamp(expr, indicatorArrayMap, timestamp)
		if err != nil {
			return nil, fmt.Errorf("error evaluating at timestamp %v: %w", timestamp, err)
		}
		results[timestamp] = result
	}

	return results, nil
}

// GetLatestConditionResult gets the latest evaluation result
func GetLatestConditionResult(
	expr string,
	indicatorArrayMap map[string][]IndicatorValues,
) (bool, time.Time, error) {
	// Find the latest timestamp
	var latestTimestamp time.Time
	var hasTimestamp bool

	for _, values := range indicatorArrayMap {
		for _, val := range values {
			if !hasTimestamp || val.Timestamp.After(latestTimestamp) {
				latestTimestamp = val.Timestamp
				hasTimestamp = true
			}
		}
	}

	if !hasTimestamp {
		return false, time.Time{}, fmt.Errorf("no data found in indicator array map")
	}

	result, err := EvaluateConditionAtTimestamp(expr, indicatorArrayMap, latestTimestamp)
	return result, latestTimestamp, err
}

// ExtractIndicators parses an expression string and returns all indicator function calls
func ExtractIndicators(expr string) ([]IndicatorInfo, error) {
	// Parse the expression into an AST
	node, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}

	// Extract all function calls from the AST
	var indicators []IndicatorInfo
	extractFunctionCalls(node, &indicators)

	return indicators, nil
}

// extractFunctionCalls recursively traverses the AST and collects all function calls
func extractFunctionCalls(node Node, indicators *[]IndicatorInfo) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *FunctionCallNode:
		// Extract arguments from the function call
		args := make([]any, len(n.Arguments))
		for i, arg := range n.Arguments {
			if varNode, isVariable := arg.(*VariableNode); isVariable {
				args[i] = varNode.Name
				continue
			}

			switch argNode := arg.(type) {
			case *NumberNode:
				args[i] = argNode.Value
			case *FunctionCallNode:
				// For nested function calls, we'll represent them as strings
				args[i] = argNode.String()
			default:
				// For other node types, convert to string representation
				args[i] = argNode.String()
			}
		}

		// Add the indicator to our list
		*indicators = append(*indicators, IndicatorInfo{
			Name:      n.Name,
			Arguments: args,
		})

		// Recursively process nested function calls in arguments
		for _, arg := range n.Arguments {
			extractFunctionCalls(arg, indicators)
		}

	case *BinaryOpNode:
		// Process both left and right operands
		extractFunctionCalls(n.Left, indicators)
		extractFunctionCalls(n.Right, indicators)

	case *UnaryOpNode:
		// Process the operand
		extractFunctionCalls(n.Operand, indicators)

	case *NumberNode, *VariableNode:
		// These nodes don't contain function calls, so we don't need to process them
		return
	}
}
