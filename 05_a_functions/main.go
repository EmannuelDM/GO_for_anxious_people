package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// 1. BASIC FUNCTION SIGNATURES & PASS-BY-VALUE
// ============================================================================

// Consecutive parameters of the same type can share a single type declaration.
func add(a, b int) int {
	return a + b
}

// Go is strictly PASS-BY-VALUE: function receives a copy of the argument.
func tryToModifyValue(num int) {
	num = 999 // Only modifies the local copy inside this function!
}

// To modify the caller's variable, pass a pointer (*int).
func modifyViaPointer(ptr *int) {
	if ptr != nil {
		*ptr = 999 // Modifies the memory address pointed to by ptr
	}
}

// ============================================================================
// 2. MULTIPLE RETURN VALUES & IDIOMATIC ERROR HANDLING
// ============================================================================

// Go functions frequently return (result, error). There are no try/catch exceptions.
func divide(numerator, denominator float64) (float64, error) {
	if denominator == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, nil
}

// ============================================================================
// 3. NAMED RETURN VALUES & NAKED RETURNS
// ============================================================================

// Named returns create variables initialized to their zero values at function start.
// 'return' without arguments returns the current values of those named variables ("naked return").
func calculateStats(numbers []int) (min int, max int, avg float64) {
	if len(numbers) == 0 {
		return // Naked return: returns (0, 0, 0.0)
	}

	min = numbers[0]
	max = numbers[0]
	sum := 0

	for _, n := range numbers {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
		sum += n
	}

	avg = float64(sum) / float64(len(numbers))
	return // Naked return: returns current min, max, avg
}

// Crucial Go Idiom: Defer functions can inspect and alter named return values before exit!
func modifyNamedReturnViaDefer() (result int) {
	defer func() {
		result *= 2 // Doubles whatever was specified in the return statement
	}()

	return 5 // result is initially 5, defer turns it into 10 before caller gets it!
}

// ============================================================================
// 4. VARIADIC FUNCTIONS (...T)
// ============================================================================

// The '...int' syntax allows zero or more integers to be passed.
// Inside the function, 'numbers' is treated as a slice: []int.
func sumNumbers(label string, numbers ...int) (string, int) {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return fmt.Sprintf("%s (count: %d)", label, len(numbers)), total
}

// ============================================================================
// 5. FIRST-CLASS & HIGHER-ORDER FUNCTIONS
// ============================================================================

// We can define custom types for function signatures for clarity and reusability.
type IntPredicate func(int) bool
type IntTransformer func(int) int

// Higher-order function: accepts a function as an argument
func filter(nums []int, isMatch IntPredicate) []int {
	var result []int
	for _, n := range nums {
		if isMatch(n) {
			result = append(result, n)
		}
	}
	return result
}

// Higher-order function: transforms a slice using a transformer function
func mapSlice(nums []int, fn IntTransformer) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

// Higher-order function: returns a new function (Function Factory)
func createMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// ============================================================================
// 6. CLOSURES & LEXICAL SCOPE
// ============================================================================

// A closure is a function that references variables from outside its body.
// The inner function retains access to 'count' even after makeCounter returns!
func makeCounter(start int, step int) func() int {
	count := start
	return func() int {
		current := count
		count += step
		return current
	}
}

// ============================================================================
// 7. DEFER MECHANICS
// ============================================================================

func demonstrateDeferStacking() {
	// Deferred calls are executed in LIFO (Last-In, First-Out) stack order.
	fmt.Println("  [defer demo] Scheduling deferred calls...")
	defer fmt.Println("  [defer demo] Execution 1: First deferred -> Executed LAST")
	defer fmt.Println("  [defer demo] Execution 2: Second deferred -> Executed MIDDLE")
	defer fmt.Println("  [defer demo] Execution 3: Third deferred -> Executed FIRST")
	fmt.Println("  [defer demo] Function body finishing, about to run deferred stack:")
}

func demonstrateDeferEvaluationTiming() {
	x := 10
	// Crucial rule: The arguments of a deferred call are evaluated IMMEDIATELY,
	// not when the function actually executes!
	defer fmt.Printf("  [defer demo] Argument evaluated at schedule time: x = %d (expected: 10)\n", x)

	// But if you use a closure, the variable is evaluated at execution time!
	defer func() {
		fmt.Printf("  [defer demo] Closure evaluated at execution time: x = %d (expected: 99)\n", x)
	}()

	x = 99
	fmt.Println("  [defer demo] Changed x to 99 before returning.")
}

// ============================================================================
// 8. PANIC AND RECOVER
// ============================================================================

// recover() can capture a panic, but ONLY when executed inside a deferred function.
func safeExecute(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert the panic value to a standard Go error
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	action()
	return nil
}

// ============================================================================
// 9. RECURSION
// ============================================================================

func factorial(n uint) uint {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	fmt.Println("==================================================")
	fmt.Println("                FUNCTIONS IN GO                   ")
	fmt.Println("==================================================")

	demoBasicSignaturesAndPointers()
	demoMultipleReturns()
	demoNamedReturns()
	demoVariadic()
	demoFirstClassFunctions()
	demoClosures()
	demoDefer()
	demoPanicAndRecover()
	demoRecursion()
}

// ----------------------------------------------------------------------------
// Section 1: Signatures and Pass-by-Value
// ----------------------------------------------------------------------------
func demoBasicSignaturesAndPointers() {
	fmt.Println("\n--- 1. Basic Signatures & Pass-by-Value ---")

	sum := add(15, 25)
	fmt.Printf("add(15, 25) = %d\n", sum)

	original := 42
	tryToModifyValue(original)
	fmt.Printf("After tryToModifyValue(42): %d (unchanged because of pass-by-value)\n", original)

	modifyViaPointer(&original)
	fmt.Printf("After modifyViaPointer(&original): %d (successfully modified)\n", original)
}

// ----------------------------------------------------------------------------
// Section 2: Multiple Returns and Errors
// ----------------------------------------------------------------------------
func demoMultipleReturns() {
	fmt.Println("\n--- 2. Multiple Returns & Error Handling ---")

	// Successful division
	res, err := divide(100, 4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("100 / 4 = %.2f\n", res)
	}

	// Division by zero
	res, err = divide(50, 0)
	if err != nil {
		fmt.Printf("50 / 0 returned error: %q (safe handling!)\n", err)
	} else {
		fmt.Printf("Result: %.2f\n", res)
	}

	// Discarding unwanted return values with blank identifier '_'
	_, errOnly := divide(10, 0)
	fmt.Printf("Only checked error using blank identifier _: %v\n", errOnly)
}

// ----------------------------------------------------------------------------
// Section 3: Named Returns and Defer Alteration
// ----------------------------------------------------------------------------
func demoNamedReturns() {
	fmt.Println("\n--- 3. Named Return Values & Naked Returns ---")

	scores := []int{45, 92, 78, 60, 88, 95, 71}
	min, max, avg := calculateStats(scores)
	fmt.Printf("Scores: %v\nMin: %d, Max: %d, Average: %.2f\n", scores, min, max, avg)

	// Modifying named return via defer
	modified := modifyNamedReturnViaDefer()
	fmt.Printf("modifyNamedReturnViaDefer() returned: %d (5 was doubled by defer!)\n", modified)
}

// ----------------------------------------------------------------------------
// Section 4: Variadic Functions
// ----------------------------------------------------------------------------
func demoVariadic() {
	fmt.Println("\n--- 4. Variadic Functions (...T) ---")

	// Calling with explicit arguments
	label, total := sumNumbers("Direct items", 10, 20, 30, 40)
	fmt.Printf("%s -> Total: %d\n", label, total)

	// Calling with zero variadic arguments
	label0, total0 := sumNumbers("Zero items")
	fmt.Printf("%s -> Total: %d\n", label0, total0)

	// Unpacking an existing slice with the spread operator '...'
	sliceItems := []int{100, 200, 300}
	labelSlice, totalSlice := sumNumbers("Unpacked slice", sliceItems...)
	fmt.Printf("%s -> Total: %d\n", labelSlice, totalSlice)
}

// ----------------------------------------------------------------------------
// Section 5: First-Class and Higher-Order Functions
// ----------------------------------------------------------------------------
func demoFirstClassFunctions() {
	fmt.Println("\n--- 5. First-Class & Higher-Order Functions ---")

	// Assigning an anonymous function to a variable
	greet := func(name string) string {
		return "Hello, " + strings.ToUpper(name) + "!"
	}
	fmt.Println(greet("Gopher"))

	// Passing functions as arguments (Higher-Order)
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	isEven := func(n int) bool { return n%2 == 0 }
	evens := filter(numbers, isEven)
	fmt.Printf("Filtered Evens: %v\n", evens)

	// Using transformer function
	square := func(n int) int { return n * n }
	squared := mapSlice(evens, square)
	fmt.Printf("Squared Evens: %v\n", squared)

	// Function factory: function returning another function
	triple := createMultiplier(3)
	fmt.Printf("createMultiplier(3) applied to 7: %d\n", triple(7))
}

// ----------------------------------------------------------------------------
// Section 6: Closures
// ----------------------------------------------------------------------------
func demoClosures() {
	fmt.Println("\n--- 6. Closures & State Encapsulation ---")

	// Each counter has its own private, encapsulated 'count' state
	counterA := makeCounter(0, 1)
	counterB := makeCounter(100, 10)

	fmt.Printf("Counter A (start 0, step 1): %d, %d, %d\n", counterA(), counterA(), counterA())
	fmt.Printf("Counter B (start 100, step 10): %d, %d, %d\n", counterB(), counterB(), counterB())
	fmt.Printf("Counter A continues independently: %d\n", counterA())
}

// ----------------------------------------------------------------------------
// Section 7: Defer Mechanics
// ----------------------------------------------------------------------------
func demoDefer() {
	fmt.Println("\n--- 7. Defer Mechanics (LIFO & Evaluation Timing) ---")
	demonstrateDeferStacking()
	fmt.Println()
	demonstrateDeferEvaluationTiming()
}

// ----------------------------------------------------------------------------
// Section 8: Panic and Recover
// ----------------------------------------------------------------------------
func demoPanicAndRecover() {
	fmt.Println("\n--- 8. Panic & Recover ---")

	// Normal operation
	err1 := safeExecute(func() {
		fmt.Println("  Executing normal routine...")
	})
	fmt.Printf("  Normal routine error: %v (nil = success)\n", err1)

	// Panicking operation
	err2 := safeExecute(func() {
		fmt.Println("  Executing dangerous routine that triggers panic...")
		var ptr *int
		*ptr = 123 // Causes runtime panic: nil pointer dereference
	})
	fmt.Printf("  Dangerous routine safely intercepted: %v\n", err2)
	fmt.Println("  Program continues running normally without crashing!")
}

// ----------------------------------------------------------------------------
// Section 9: Recursion
// ----------------------------------------------------------------------------
func demoRecursion() {
	fmt.Println("\n--- 9. Recursion ---")

	for n := uint(0); n <= 6; n++ {
		fmt.Printf("  factorial(%d) = %d\n", n, factorial(n))
	}
}
