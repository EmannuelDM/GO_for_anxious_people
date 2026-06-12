package main

import "fmt"

// Person struct is used to demonstrate pointers with structs
type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("               POINTERS IN GO                     ")
	fmt.Println("==================================================")

	// 1. Declaring and Initializing Pointers
	// A pointer holds the memory address of a value.
	// The type *T is a pointer to a T value.
	var x int = 42
	var p *int // Declares a pointer to an integer. Its zero value is nil.

	fmt.Println("\n--- 1. Declaration and Initialization ---")
	fmt.Printf("x (value): %v, type: %T\n", x, x)
	fmt.Printf("p (initial zero value): %v, type: %T\n", p, p)

	// Use the & operator to get the address of a variable
	p = &x
	fmt.Printf("p (holding address of x): %v\n", p)
	fmt.Printf("Address of x (&x): %v\n", &x)

	// 2. Dereferencing (Accessing the value)
	// Use the * operator to access/read the value at the address (dereferencing)
	fmt.Println("\n--- 2. Dereferencing (*) ---")
	fmt.Printf("Value pointed to by p (*p): %v\n", *p)

	// Modifying the value through the pointer
	*p = 100
	fmt.Println("Modifying value via pointer: *p = 100")
	fmt.Printf("New value of x: %v\n", x)
	fmt.Printf("New value of *p: %v\n", *p)

	// 3. Multiple Pointers
	// Multiple pointers can point to the same memory address
	var q *int = p
	fmt.Println("\n--- 3. Multiple Pointers ---")
	fmt.Printf("q (points to address in p): %v, value (*q): %d\n", q, *q)
	*q = 500
	fmt.Println("Modifying value via q: *q = 500")
	fmt.Printf("Value of x: %d, *p: %d, *q: %d\n", x, *p, *q)

	// 4. Passing Pointers to Functions
	// Go passes arguments by value (making a copy of the argument).
	// Passing a pointer allows a function to modify the original variable.
	fmt.Println("\n--- 4. Pointer Functions vs Value Functions ---")
	val := 10
	fmt.Printf("Initial val: %d\n", val)

	modifyValue(val)
	fmt.Printf("After modifyValue(val): %d (No change, passed by value)\n", val)

	modifyPointer(&val)
	fmt.Printf("After modifyPointer(&val): %d (Value changed, passed by pointer)\n", val)

	// 5. The new() Built-in Function
	// The new(T) function allocates zeroed memory for type T and returns its address (*T).
	fmt.Println("\n--- 5. The new() Function ---")
	newPtr := new(int) // allocates memory for an int, sets its value to 0, returns *int
	fmt.Printf("newPtr: address = %v, value (*newPtr) = %d\n", newPtr, *newPtr)
	*newPtr = 777
	fmt.Printf("newPtr after modification: address = %v, value (*newPtr) = %d\n", newPtr, *newPtr)

	// 6. Pointers to Structs
	// Structs are very common in Go. Let's see how pointers interact with them.
	fmt.Println("\n--- 6. Pointers to Structs ---")
	person := Person{Name: "Alice", Age: 30}
	fmt.Printf("Original struct: %+v\n", person)

	// Obtain pointer to struct
	personPtr := &person

	// Syntactic Sugar: In Go, we can write personPtr.Age instead of (*personPtr).Age.
	// The language automatically dereferences the struct pointer for us.
	personPtr.Age = 31
	personPtr.Name = "Bob"
	fmt.Printf("After modifying via pointer (personPtr.Age = 31, Name = Bob): %+v\n", person)

	// 7. No Pointer Arithmetic
	// Unlike C/C++, Go does not allow arithmetic on pointers (e.g., p++ or p + 2) by default.
	// This prevents bugs related to invalid memory access.
	fmt.Println("\n--- 7. Pointer Arithmetic Safety ---")
	fmt.Println("Go does not allow operations like p++ or p + 1. The following would cause a compilation error:")
	fmt.Println("    p = p + 1")
}

// modifyValue takes an integer by value (creates a local copy)
func modifyValue(v int) {
	v = 999
}

// modifyPointer takes a pointer to an integer (can modify the original value)
func modifyPointer(v *int) {
	*v = 999
}
