package main

import "fmt"

func main() {
	fmt.Println("==================================================")
	fmt.Println("                ARRAYS IN GO                      ")
	fmt.Println("==================================================")

	// 1. Declaration and Zero Values
	// An array is a numbered sequence of elements of a single type with a fixed length.
	// The length is part of the type, so [5]int and [10]int are different types.
	// Arrays are zero-valued by default (elements get their default zero value).
	fmt.Println("\n--- 1. Declaration and Zero Values ---")
	var a [5]int // declares an array of 5 integers, initialized to 0
	fmt.Printf("var a [5]int -> Value: %v, Type: %T, Length: %d\n", a, a, len(a))

	var strings [3]string // zero value of string is ""
	fmt.Printf("var strings [3]string -> Value: %q, Type: %T\n", strings, strings)

	// 2. Initialization with Array Literals
	fmt.Println("\n--- 2. Initialization with Literals ---")
	// Initialize with specific values
	b := [5]int{10, 20, 30, 40, 50}
	fmt.Printf("b := [5]int{10, 20, 30, 40, 50} -> Value: %v\n", b)

	// Sparse initialization (specific indices)
	// Here, index 1 is 12, index 4 is 42, and others are 0
	c := [5]int{1: 12, 4: 42}
	fmt.Printf("c := [5]int{1: 12, 4: 42} -> Value: %v\n", c)

	// 3. Compiler-counted Length
	fmt.Println("\n--- 3. Let the Compiler Count Length (...) ---")
	// You can use '...' to let the compiler determine the size from the number of elements
	d := [...]int{1, 2, 3, 5, 8, 13}
	fmt.Printf("d := [...]int{1, 2, 3, 5, 8, 13} -> Value: %v, Type: %T, Length: %d\n", d, d, len(d))

	// 4. Accessing and Modifying Elements
	fmt.Println("\n--- 4. Accessing and Modifying ---")
	var colors [3]string
	colors[0] = "Red"
	colors[1] = "Green"
	colors[2] = "Blue"
	fmt.Printf("colors after modification: %v\n", colors)
	fmt.Printf("colors[1]: %s\n", colors[1])

	// 5. Arrays are Value Types
	// In Go, arrays are value types, not reference types.
	// Assigning an array to another, or passing it to a function, creates a COMPLETE COPY of its data.
	fmt.Println("\n--- 5. Arrays are Value Types (Passed/Assigned by Copy) ---")
	original := [3]int{1, 2, 3}
	copyAr := original // makes a complete copy

	copyAr[0] = 999
	fmt.Printf("original array: %v (Unchanged!)\n", original)
	fmt.Printf("copied array  : %v\n", copyAr)

	// Let's test passing an array to functions
	modifyByValue(original)
	fmt.Printf("After modifyByValue(original): %v (No change! passed by value)\n", original)

	modifyByPointer(&original)
	fmt.Printf("After modifyByPointer(&original): %v (Changed! passed address)\n", original)

	// 6. Multi-dimensional Arrays
	fmt.Println("\n--- 6. Multi-dimensional Arrays ---")
	// Array types are one-dimensional, but you can compose types to build multi-dimensional arrays.
	var matrix [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			matrix[i][j] = i + j
		}
	}
	fmt.Printf("matrix [2][3]int: %v\n", matrix)

	// Literal initialization of a 2D array
	grid := [2][2]string{
		{"A1", "A2"},
		{"B1", "B2"},
	}
	fmt.Printf("grid: %v\n", grid)

	// 7. Iterating Over Arrays
	fmt.Println("\n--- 7. Iteration (Looping) ---")
	fruits := [...]string{"Apple", "Banana", "Cherry"}

	// Using standard for loop
	fmt.Println("Using standard for loop:")
	for i := 0; i < len(fruits); i++ {
		fmt.Printf("  Index %d: %s\n", i, fruits[i])
	}

	// Using range loop (index and value)
	fmt.Println("Using range loop (index and value):")
	for index, value := range fruits {
		fmt.Printf("  Index %d: %s\n", index, value)
	}

	// Using range loop (value only, ignoring index)
	fmt.Println("Using range loop (ignoring index with _):")
	for _, value := range fruits {
		fmt.Printf("  Value: %s\n", value)
	}

	// 8. Array Comparison
	// You can compare arrays of the same type (same length and element type) using == and !=.
	// Arrays are equal if their corresponding elements are equal.
	fmt.Println("\n--- 8. Array Comparison ---")
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	arr3 := [3]int{1, 2, 4}
	// Note: We cannot compare [3]int with [4]int; it would be a compilation error.
	fmt.Printf("arr1 == arr2: %t\n", arr1 == arr2)
	fmt.Printf("arr1 == arr3: %t\n", arr1 == arr3)
}

// modifyByValue receives a copy of the array
func modifyByValue(arr [3]int) {
	arr[0] = 555
}

// modifyByPointer receives a pointer to the array, allowing modifications
func modifyByPointer(arr *[3]int) {
	arr[0] = 555
}
