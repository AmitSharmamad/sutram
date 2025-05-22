# Expression Parser Module Documentation

This document provides comprehensive documentation on the Expression Parser module's grammar, syntax, and available functions. This can be used to build any valid expression with the syntax understood by the parser module.

## 1. Introduction

The Expression Parser module allows users to write mathematical and logical expressions that can be parsed and evaluated in the context of trading strategies. It provides a flexible way to express entry and exit conditions, technical indicators, and other trading rules.

### 1.1 Expression Processing

The parser takes string expressions written in the defined syntax and converts them into an Abstract Syntax Tree (AST) which can then be evaluated against market data. This allows for dynamic strategy criteria without requiring code changes.

## 2. Expression Grammar

The parser uses a formal grammar with the following components:

### 2.1 Basic Elements

- **Numbers**: Decimal literals (e.g., `42`, `3.14`)
- **Variables**: Named references to data (e.g., `CLOSE`, `VOLUME`)
- **Operators**: Mathematical, comparison, and logical operations
- **Functions**: Named operations with parameters (e.g., `SMA(CLOSE, 20)`)
- **Parentheses**: For grouping expressions and function arguments

### 2.2 Expression Structure

Expressions follow this fundamental structure:
```
(<Indicator> <Operator> <Indicator>)
```

Where:
- `<Indicator>` can be a variable, number, function call, or another nested expression
- `<Operator>` can be an arithmetic, comparison, or logical operator

Expressions can be nested to create complex conditions:
```
(<Indicator1> <Operator1> <Indicator2>) <LogicalOperator> (<Indicator3> <Operator2> <Indicator4>)
```

### 2.3 Grammar Rules

The expression grammar follows a hierarchical structure with the following precedence (from highest to lowest):

1. Literals (numbers and variables) and parenthesized expressions
2. Function calls
3. Unary operations (`+`, `-`, `NOT`)
4. Multiplicative operations (`*`, `/`, `%`)
5. Additive operations (`+`, `-`)
6. Comparison operations (`>`, `<`, `>=`, `<=`, `==`, `!=`)
7. Logical AND operation
8. Logical OR operation

## 3. Syntax Elements

### 3.1 Numbers

Numbers can be integers or floating-point values:
```
42
3.14
0.5
```

### 3.2 Variables

Variables are case-insensitive identifiers representing market data:
```
close    // Current closing price
open     // Current opening price
high     // Current high price
low      // Current low price
volume   // Current volume
```

**Important Note**: Although the parser internally processes variables in uppercase (e.g., `CLOSE`, `OPEN`), in your expressions you should use lowercase (e.g., `close`, `open`) as shown in the examples.

### 3.3 Operators

#### 3.3.1 Arithmetic Operators
- `+` : Addition
- `-` : Subtraction
- `*` : Multiplication
- `/` : Division
- `%` : Modulo (remainder)

#### 3.3.2 Comparison Operators
- `>` : Greater than
- `<` : Less than
- `>=` : Greater than or equal to
- `<=` : Less than or equal to
- `==` : Equal to
- `!=` : Not equal to

#### 3.3.3 Logical Operators
- `AND` : Logical AND
- `OR` : Logical OR
- `NOT` : Logical NOT

### 3.4 Parentheses

Parentheses `()` are used to:
1. Override operator precedence
2. Group expressions
3. Define function arguments

## 4. Technical Indicators and Functions

The parser supports various technical indicators and functions that can be used as part of expressions:

### 4.1 Price Data

These are used directly as variables in expressions:
```
close     // Current closing price
open      // Current opening price
high      // Current high price
low       // Current low price
volume    // Current volume
```

### 4.2 Moving Averages

- `sma(source, period)` : Simple Moving Average
  Example: `sma(close, 20)` - 20-period Simple Moving Average of closing prices
  
- `ema(source, period, offset)` : Exponential Moving Average
  Example: `ema(close, 9, 0)` - 9-period Exponential Moving Average of closing prices

### 4.3 Oscillators and Momentum Indicators

- `rsi(source, period)` : Relative Strength Index
  Example: `rsi(close, 14)` - 14-period RSI of closing prices
  
- `macd(source, fastPeriod, slowPeriod)` : MACD Line
  Example: `macd(close, 12, 26)` - MACD line using 12 and 26 periods
  
- `macd_signal(source, fastPeriod, slowPeriod, signalPeriod)` : MACD Signal Line
  Example: `macd_signal(close, 12, 26, 9)` - 9-period signal line of MACD
  
- `stoch_k(kPeriod, kSlowing)` : Stochastic %K
  Example: `stoch_k(14, 3)` - Stochastic %K with period 14 and slowing 3
  
- `stoch_d(kPeriod, kSlowing)` : Stochastic %D
  Example: `stoch_d(14, 3)` - Stochastic %D with period 14 and slowing 3
  
- `tsi(source, longPeriod, shortPeriod)` : True Strength Index
  Example: `tsi(close, 25, 13)` - TSI using 25 and 13 periods
  
- `mfi(period)` : Money Flow Index
  Example: `mfi(14)` - 14-period Money Flow Index
  
- `uo(period1, period2, period3)` : Ultimate Oscillator
  Example: `uo(7, 14, 28)` - Ultimate Oscillator with 7, 14, and 28 periods
  
- `vo(shortPeriod, longPeriod)` : Volume Oscillator
  Example: `vo(5, 10)` - Volume Oscillator with 5 and 10 periods

### 4.4 Volatility Indicators

- `bollinger(source, period, deviation)` : Bollinger Bands Middle Band
  Example: `bollinger(close, 20, 2)` - 20-period Bollinger Band middle band with 2 standard deviations
  
- `upper_band(source, period, deviation)` : Bollinger Bands Upper Band
  Example: `upper_band(close, 20, 2)` - Upper band of Bollinger Bands
  
- `lower_band(source, period, deviation)` : Bollinger Bands Lower Band
  Example: `lower_band(close, 20, 2)` - Lower band of Bollinger Bands
  
- `atr(period)` : Average True Range
  Example: `atr(14)` - 14-period Average True Range
  
- `supertrend(period, multiplier)` : SuperTrend Indicator
  Example: `supertrend(10, 3)` - SuperTrend with period 10 and multiplier 3
  
- `keltner(source, period, atrMultiplier)` : Keltner Channels Middle Line
  Example: `keltner(close, 20, 2)` - 20-period Keltner Channel middle line with ATR multiplier 2

### 4.5 Trend Indicators

- `dmi_plus(period)` : Positive Directional Movement Index (+DI)
  Example: `dmi_plus(14)` - 14-period +DI
  
- `dmi_minus(period)` : Negative Directional Movement Index (-DI)
  Example: `dmi_minus(14)` - 14-period -DI
  
- `adx(period)` : Average Directional Index
  Example: `adx(14)` - 14-period ADX
  
- `vortex_plus(period)` : Vortex Indicator Positive
  Example: `vortex_plus(14)` - 14-period positive Vortex indicator
  
- `vortex_minus(period)` : Vortex Indicator Negative
  Example: `vortex_minus(14)` - 14-period negative Vortex indicator
  
- `psar(acceleration, maximum)` : Parabolic SAR
  Example: `psar(0.02, 0.2)` - Parabolic SAR with acceleration 0.02 and maximum 0.2
  
- `ichimoku_tenkan(period)` : Ichimoku Tenkan-sen (Conversion Line)
  Example: `ichimoku_tenkan(9)` - 9-period Tenkan-sen
  
- `ichimoku_kijun(period)` : Ichimoku Kijun-sen (Base Line)
  Example: `ichimoku_kijun(26)` - 26-period Kijun-sen
  
- `ichimoku_cloud_top(tenkanPeriod, kijunPeriod, spanBPeriod)` : Ichimoku Cloud Top (Senkou Span A)
  Example: `ichimoku_cloud_top(9, 26, 52)` - Senkou Span A with standard periods
  
- `ichimoku_cloud_bottom(tenkanPeriod, kijunPeriod, spanBPeriod)` : Ichimoku Cloud Bottom (Senkou Span B)
  Example: `ichimoku_cloud_bottom(9, 26, 52)` - Senkou Span B with standard periods
  
- `vwap()` : Volume Weighted Average Price
  Example: `vwap()` - VWAP for the current session

### 4.6 Utility Functions

- `crossover(indicator1, indicator2)` : Returns true when indicator1 crosses above indicator2
  Example: `crossover(ema(close, 9), ema(close, 21))` - True when 9-period EMA crosses above 21-period EMA
  
- `crossunder(indicator1, indicator2)` : Returns true when indicator1 crosses below indicator2
  Example: `crossunder(ema(close, 9), ema(close, 21))` - True when 9-period EMA crosses below 21-period EMA

**Important Note**: While the SQL examples use `crossover` and `crossunder` functions, the parser implementation internally uses `CROSSESABOVE` and `CROSSESBELOW` operators. The system automatically maps between these, so you should use `crossover` and `crossunder` in your expressions.
  
- `highest(source, period)` : Highest value over a given period
  Example: `highest(high, 20)` - Highest high over the last 20 periods
  
- `lowest(source, period)` : Lowest value over a given period
  Example: `lowest(low, 20)` - Lowest low over the last 20 periods
  
- `roc(source, period)` : Rate of Change
  Example: `roc(close, 10)` - 10-period rate of change of closing prices
  
- `range(period)` : Range (high-low) over a given period
  Example: `range(5)` - Average range over the last 5 periods
  
- `high_of_last(period)` : Highest high of the last n periods
  Example: `high_of_last(20)` - Highest high over the last 20 periods
  
- `low_of_last(period)` : Lowest low of the last n periods
  Example: `low_of_last(20)` - Lowest low over the last 20 periods
  
- `trend(source, period)` : Trend direction (positive or negative)
  Example: `trend(close, 10)` - Trend direction of the last 10 closing prices

## 5. Expression Evaluation

Expressions are evaluated to produce numeric results:
- Arithmetic operations produce numeric values
- Comparison operations produce 1 (true) or 0 (false)
- Logical operations produce 1 (true) or 0 (false)

The result of evaluating the entire expression is interpreted as follows:
- Values greater than 0 are considered TRUE
- Values equal to 0 are considered FALSE

## 6. Examples

### 6.1 Simple Expressions

```
close > open                          // True when close price is higher than open price
volume > 1000000                      // True when volume exceeds 1 million
(high - low) > 2 * atr(14)            // True when the day's range exceeds 2x ATR
```

### 6.2 Moving Average Strategies

```
close > sma(close, 20)                // True when price is above 20-day SMA
crossover(close, sma(close, 50))      // True when price crosses above 50-day SMA
ema(close, 9, 0) > ema(close, 21, 0)  // True when 9-day EMA is above 21-day EMA
```

### 6.3 Technical Indicator Combinations

```
rsi(close, 14) < 30                                      // Oversold condition
macd(close, 12, 26) > 0 AND rsi(close, 14) > 50         // MACD positive and RSI bullish
supertrend(10, 3) < close AND rsi(close, 14) > 50       // SuperTrend bullish and RSI bullish
```

### 6.4 Complex Strategies

```
(crossover(close, upper_band(close, 20, 2)) OR rsi(close, 14) < 30) AND volume > sma(volume, 20)
```
This complex expression evaluates to true when:
- Either the price crosses above the upper Bollinger Band or RSI is below 30 (oversold)
- AND the volume is above its 20-day average

### 6.5 Nested Expressions

```
(ema(close, 9, 0) > ema(close, 21, 0)) AND ((rsi(close, 14) > 50) OR (volume > sma(volume, 20) * 1.5))
```
This nested expression evaluates to true when:
- The 9-day EMA is above the 21-day EMA
- AND either RSI is above 50 OR volume is 50% above its 20-day average

### 6.6 Real Strategy Examples

Here are examples of complete strategy expressions used in actual trading:

```
// Moving Average Crossover Strategy
close > ema(close, 9) AND ema(close, 9) > ema(close, 21) AND crossover(ema(close, 9), ema(close, 21))

// RSI Reversal Strategy
rsi(close, 14) < 30 AND rsi(close, 14)[1] < rsi(close, 14)

// Bollinger Band Bounce Strategy
low <= lower_band(close, 20, 2) AND close > open AND close[1] < open[1]

// MACD Bullish Crossover Strategy
crossover(macd(close, 12, 26), macd_signal(close, 12, 26, 9)) AND macd(close, 12, 26) < 0

// Stochastic Crossover Strategy
stoch_k(14, 3) < 30 AND stoch_d(14, 3) < 30 AND crossover(stoch_k(14, 3), stoch_d(14, 3))
```

## 7. Best Practices

1. **Use Parentheses**: Use parentheses to clearly indicate the order of operations and group related expressions, especially in complex nested expressions.

2. **Proper Syntax**: Always follow the `(<Indicator> <Operator> <Indicator>)` structure for expressions, ensuring each expression is well-formed.

3. **Case Sensitivity**: While the parser is case-insensitive, it's best to use lowercase for variables and function names as shown in the examples (e.g., `close`, `ema`, `crossover`).

4. **Avoid Division by Zero**: The parser will handle division by zero by returning an error, but it's best to avoid such situations in your expressions.

5. **Historical Lookback**: Be aware that some indicators require historical data to calculate. Ensure your data source provides sufficient history.

6. **Function Names vs. Operators**: Remember that `crossover` and `crossunder` are implemented as functions in the examples but are processed as operators (`CROSSESABOVE` and `CROSSESBELOW`) in the parser.

7. **Test Thoroughly**: Test your expressions with historical data before using them in live trading.

8. **Readability**: Break complex expressions into smaller parts for better readability and maintenance.

## 8. Limitations

1. **Syntax Errors**: The parser will return errors for invalid syntax or undefined variables/functions.

2. **Data Availability**: Expressions requiring historical data beyond what's available will return errors.

3. **Recursive Function Calls**: The parser does not support recursive function calls.

4. **Valid Indicator-Operator Combinations**: Ensure that you use operators that make sense for the indicators they connect. For example, comparing a momentum oscillator to a price might not be meaningful.

5. **Standard Indicators Only**: Only the indicators and functions documented here are available. Custom indicators must be implemented separately.

6. **Internal vs. External Representation**: There are some differences between how expressions are represented in the SQL examples and how they're processed internally by the parser. Use the function and operator names as shown in the examples.

## 9. LLM Integration Guidelines

When using an LLM to generate expressions for the parser module:

1. **Explain the Syntax**: Clearly explain the `(<Indicator> <Operator> <Indicator>)` syntax and the available operators.

2. **Provide Function Reference**: Share the comprehensive list of available indicators and their parameters.

3. **Include Examples**: Show examples of valid expressions, especially ones that demonstrate nesting and logical combinations.

4. **Validate Output**: Have the LLM validate its own expressions for syntax correctness.

5. **Request Explanations**: Have the LLM explain the logic of generated expressions to ensure they match the intended trading strategy.

## 10. Summary

The Expression Parser module provides a powerful way to express trading conditions using a structured syntax of `(<Indicator> <Operator> <Indicator>)` with support for nesting and logical operators. By combining price data, technical indicators, and logical operations, users can create sophisticated trading strategies that can be automatically evaluated against market data.

### 10.1 Implementation Details

The parser internally transforms expressions into an Abstract Syntax Tree (AST) with different node types:
- NumberNode: Represents numeric literals
- VariableNode: Represents variables like `close`, `open`
- BinaryOpNode: Represents operations like `>`, `AND`, `CROSSESABOVE`
- UnaryOpNode: Represents unary operations like `-`, `NOT`
- FunctionCallNode: Represents function calls like `ema(close, 9)`

This structured representation allows for efficient evaluation of expressions against market data.