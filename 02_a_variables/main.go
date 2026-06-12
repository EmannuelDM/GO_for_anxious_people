package main

import "fmt"

func main() {
	// 1. Boolean type
	var isGoFun bool = true

	// 2. String type
	var message string = "Hello, Go!"

	// 3. Integer types (signed and unsigned)
	var age int = 30
	var negativeInt int8 = -120
	var unsignedInt uint16 = 65535

	// 4. Floating point types
	var price float32 = 19.99
	var pi float64 = 3.1415926535

	// 5. Complex types
	var complexNum complex128 = 1 + 4i

	// 6. Type Inference (compiler decides the type)
	var inferredString = "Inferred Type"
	var inferredInt = 42

	// 7. Short Variable Declaration (only inside functions)
	shortVar := 99.9

	// 8. Zero Values (default values when declared without initialization)
	var defaultInt int
	var defaultBool bool
	var defaultString string
	var defaultFloat float64

	// Print helper function to display value and type
	printTypeAndValue := func(name string, val interface{}) {
		fmt.Printf("%-18s | Value: %-15v | Type: %T\n", name, val, val)
	}

	fmt.Println("=== Declared Variables ===")
	printTypeAndValue("isGoFun", isGoFun)
	printTypeAndValue("message", message)
	printTypeAndValue("age", age)
	printTypeAndValue("negativeInt", negativeInt)
	printTypeAndValue("unsignedInt", unsignedInt)
	printTypeAndValue("price", price)
	printTypeAndValue("pi", pi)
	printTypeAndValue("complexNum", complexNum)
	printTypeAndValue("inferredString", inferredString)
	printTypeAndValue("inferredInt", inferredInt)
	printTypeAndValue("shortVar", shortVar)

	fmt.Println("\n=== Zero Values (Defaults) ===")
	printTypeAndValue("defaultInt", defaultInt)
	printTypeAndValue("defaultBool", defaultBool)
	printTypeAndValue("defaultString", fmt.Sprintf("%q", defaultString)) // Quote to show empty string
	printTypeAndValue("defaultFloat", defaultFloat)
}
