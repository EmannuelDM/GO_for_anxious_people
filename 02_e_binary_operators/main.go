package main

import "fmt"

func main() {
	fmt.Println("==================================================")
	fmt.Println("            BINARY OPERATORS IN GO                ")
	fmt.Println("==================================================")

	// 1. Arithmetic Operators
	// +, -, *, /, %
	fmt.Println("\n--- 1. Arithmetic Operators ---")
	a, b := 10, 3
	fmt.Printf("%d + %d = %d\n", a, b, a+b)
	fmt.Printf("%d - %d = %d\n", a, b, a-b)
	fmt.Printf("%d * %d = %d\n", a, b, a*b)

	// Division: Integer division truncates towards zero
	fmt.Printf("%d / %d = %d (Integer division truncates)\n", a, b, a/b)
	// For fractional division, at least one operand must be a float type
	fmt.Printf("%.1f / %.1f = %.2f (Float division)\n", float64(a), float64(b), float64(a)/float64(b))

	// Modulo: Remainder of integer division (only valid for integers)
	fmt.Printf("%d %% %d = %d (Modulo/Remainder)\n", a, b, a%b)

	// 2. Comparison Operators
	// ==, !=, <, <=, >, >=
	// Comparison operators yield a boolean value. In Go, both operands must be of the same type.
	fmt.Println("\n--- 2. Comparison Operators ---")
	fmt.Printf("%d == %d is %t\n", a, b, a == b)
	fmt.Printf("%d != %d is %t\n", a, b, a != b)
	fmt.Printf("%d > %d  is %t\n", a, b, a > b)
	fmt.Printf("%d < %d  is %t\n", a, b, a < b)

	// 3. Logical Operators
	// && (AND), || (OR)
	// These only apply to boolean values and support short-circuit evaluation.
	fmt.Println("\n--- 3. Logical Operators ---")
	t, f := true, false
	fmt.Printf("true && false is %t\n", t && f)
	fmt.Printf("true || false is %t\n", t || f)

	// Demonstrating short-circuiting:
	// In (false && printAndReturnTrue()), the function is never called because the outcome is already known.
	fmt.Println("Short-circuit evaluation check:")
	_ = f && printAndReturnTrue("AND") // Nothing printed
	_ = t || printAndReturnTrue("OR")  // Nothing printed

	// 4. Bitwise Operators
	// & (AND), | (OR), ^ (XOR), &^ (Bit Clear / AND NOT), << (Left Shift), >> (Right Shift)
	fmt.Println("\n--- 4. Bitwise Operators ---")
	// Let's use 10 (binary 1010) and 3 (binary 0011)
	x, y := 10, 3 // x = 1010, y = 0011 in binary
	fmt.Printf("x (10) in binary: %04b\n", x)
	fmt.Printf("y (3)  in binary: %04b\n", y)

	fmt.Printf("x & y (Bitwise AND):  %04b (Decimal %d)\n", x&y, x&y)   // 1010 & 0011 = 0010 (2)
	fmt.Printf("x | y (Bitwise OR):   %04b (Decimal %d)\n", x|y, x|y)   // 1010 | 0011 = 1011 (11)
	fmt.Printf("x ^ y (Bitwise XOR):  %04b (Decimal %d)\n", x^y, x^y)   // 1010 ^ 0011 = 1001 (9)

	// &^ Bit Clear (AND NOT): clears bits in the first operand where the second operand has 1s.
	// Formula: x &^ y is equivalent to x & (^y)
	fmt.Printf("x &^ y (Bit Clear):   %04b (Decimal %d)\n", x&^y, x&^y) // 1010 &^ 0011 = 1000 (8)

	// Shift Operators (left shift <<, right shift >>)
	fmt.Printf("y << 1 (Left Shift):  %04b (Decimal %d) - shifts bits left, fills with 0 (multiplies by 2)\n", y<<1, y<<1)
	fmt.Printf("x >> 1 (Right Shift): %04b (Decimal %d) - shifts bits right (integer division by 2)\n", x>>1, x>>1)

	// 5. Operator Precedence
	// Go has 5 precedence levels for binary operators (from 5=highest to 1=lowest):
	// 5:  *  /  %  <<  >>  &  &^
	// 4:  +  -  |  ^
	// 3:  ==  !=  <  <=  >  >=
	// 2:  &&
	// 1:  ||
	fmt.Println("\n--- 5. Operator Precedence ---")

	// Example 1: Multiplicative (*) has higher precedence than Additive (+)
	// 5 + 2 * 3 -> 5 + (2 * 3) = 11
	fmt.Printf("5 + 2 * 3 = %d\n", 5+2*3)

	// Example 2: Bitwise shift (<<) has higher precedence than addition (+) in Go!
	// Note: This differs from C/C++ where shift has LOWER precedence than addition.
	// In Go:  1 << 2 + 1  is evaluated as  (1 << 2) + 1  = 4 + 1 = 5
	// In C:   1 << 2 + 1  is evaluated as  1 << (2 + 1)  = 1 << 3 = 8
	fmt.Printf("1 << 2 + 1 = %d (evaluated as (1 << 2) + 1)\n", 1<<2+1)

	// Example 3: Bitwise AND (&) has higher precedence than comparison (==) in Go!
	// Note: This also differs from C/C++ where comparison has higher precedence than bitwise AND.
	// In Go:  x & y == 2  is evaluated as  (x & y) == 2 -> 2 == 2 -> true
	// In C:   x & y == 2  is evaluated as  x & (y == 2) -> 10 & false -> 0 (false)
	fmt.Printf("x & y == 2 = %t (evaluated as (x & y) == 2)\n", x&y == 2)

	// Using parentheses () to force evaluation order
	fmt.Printf("Forced with (): 1 << (2 + 1) = %d\n", 1<<(2+1))
	fmt.Printf("Forced with (): x & (y == 2) -> compilation error in Go (cannot compare int with bool)\n")
}

// printAndReturnTrue is a helper function to demonstrate logical short-circuiting
func printAndReturnTrue(opName string) bool {
	fmt.Printf("   [Side Effect] Evaluated right-hand side of %s!\n", opName)
	return true
}
