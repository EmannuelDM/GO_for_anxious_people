package mathutils

import "fmt"

// init is a special function that runs automatically when the package is initialized.
// You can have multiple init functions per package or even file.
func init() {
	fmt.Println("[mathutils/arithmetic.go] Package mathutils (arithmetic) initialized!")
}

// Add returns the sum of two integers.
// Since it starts with an uppercase letter, it is exported and accessible outside the package.
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two integers.
// Exported.
func Subtract(a, b int) int {
	return a - b
}

// multiply is a private helper function.
// Since it starts with a lowercase letter, it is unexported and can only be used
// within the mathutils package (i.e. in arithmetic.go or geometry.go).
func multiply(a, b int) int {
	return a * b
}
