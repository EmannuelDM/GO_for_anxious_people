package main

import "fmt"

func main() {
	// --- Go Keywords (25 total) ---
	// Keywords are reserved words that cannot be used as identifiers (like variable names).
	keywords := []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}

	// --- Predeclared Constants ---
	predeclaredConstants := []string{
		"true", "false", "iota", "nil",
	}

	// --- Predeclared Types ---
	predeclaredTypes := []string{
		"bool", "byte", "complex64", "complex128", "error",
		"float32", "float64", "int", "int8", "int16", "int32",
		"int64", "rune", "string", "uint", "uint8", "uint16",
		"uint32", "uint64", "uintptr",
	}

	// --- Predeclared Functions ---
	predeclaredFunctions := []string{
		"append", "cap", "close", "complex", "copy", "delete",
		"imag", "len", "make", "new", "panic", "print", "println",
		"real", "recover",
	}

	// 1. Displaying Keywords, Constants, Types, and Functions
	fmt.Println("==================================================")
	fmt.Println("           GO PREDECLARED IDENTIFIERS             ")
	fmt.Println("==================================================")

	fmt.Println("\n--- 25 Reserved Keywords ---")
	printGrid(keywords, 5)

	fmt.Println("\n--- Predeclared Constants ---")
	printGrid(predeclaredConstants, 4)

	fmt.Println("\n--- Predeclared Types ---")
	printGrid(predeclaredTypes, 5)

	fmt.Println("\n--- Predeclared Functions ---")
	printGrid(predeclaredFunctions, 5)

	fmt.Println("\n==================================================")
	fmt.Println("             DEMONSTRATING PREDECLARED            ")
	fmt.Println("==================================================")

	// Demonstrating Constants: true, false, nil
	var isTrue bool = true
	var isFalse bool = false
	var ptr *int = nil

	fmt.Printf("Constants: true=%v, false=%v, nil=%v\n", isTrue, isFalse, ptr)

	// Demonstrating Predeclared Functions: make, append, len, cap
	// slice creation using make
	slice := make([]int, 0, 5)
	fmt.Printf("After make([]int, 0, 5): len=%d, cap=%d, value=%v\n", len(slice), cap(slice), slice)

	// slice modification using append
	slice = append(slice, 10, 20, 30)
	fmt.Printf("After append(slice, 10, 20, 30): len=%d, cap=%d, value=%v\n", len(slice), cap(slice), slice)

	// Demonstrating Predeclared Functions: copy
	dest := make([]int, len(slice))
	copiedElements := copy(dest, slice)
	fmt.Printf("After copy(dest, slice): copied=%d, dest=%v\n", copiedElements, dest)

	// Demonstrating Predeclared Functions: delete (map)
	m := make(map[string]int)
	m["apple"] = 5
	m["banana"] = 10
	fmt.Printf("Map before delete: %v\n", m)
	delete(m, "apple")
	fmt.Printf("Map after delete('apple'): %v\n", m)

	// Demonstrating Predeclared Functions: complex, real, imag
	c := complex(3.0, 4.0) // 3 + 4i
	fmt.Printf("Complex number: %v (Real part: %g, Imaginary part: %g)\n", c, real(c), imag(c))
}

// printGrid prints a slice of strings in a formatted grid with `cols` columns
func printGrid(items []string, cols int) {
	for i, item := range items {
		fmt.Printf("%-12s", item)
		if (i+1)%cols == 0 {
			fmt.Println()
		}
	}
	if len(items)%cols != 0 {
		fmt.Println()
	}
}
