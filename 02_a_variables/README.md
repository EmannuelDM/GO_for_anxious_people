# Go Variables (02_a_variables)

This example demonstrates how to declare different types of variables in Go and inspect their values and types.

## Topics Covered
1. **Explicit Type Declaration**: `var name type = value`
2. **Type Inference**: `var name = value`
3. **Short Variable Declaration**: `name := value` (only works inside functions)
4. **Zero Values**: The default values assigned to variables declared without an initial value (e.g. `0` for numbers, `false` for booleans, `""` for strings).
5. **Type Inspection**: Printing variable types and values using `fmt.Printf` with the `%T` and `%v` verbs.

## fmt.Printf
Printf has over a dozen such conversions, which Go programmers call verbs.

- `%d` decimal integer
- `%x`, `%o`, `%b` integer in hexadecimal, octal, binary
- `%f`, `%g`, `%e` floating-point number: 3.141593 3.141592653589793 3.141593e+00
- `%t` boolean: true or false
- `%c` rune (Unicode code point)
- `%s` string
- `%q` quoted string "abc" or rune 'c'
- `%v` any value in a natural format
- `%T` type of any value
- `%%` literal percent sign (no operand)

## Running the Program

Navigate to this directory (`02_a_variables`) and run:

```bash
go run main.go
```
