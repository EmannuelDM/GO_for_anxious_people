package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"06_c_interfaces/service"
)

func main() {
	fmt.Println("##################################################################")
	fmt.Println("                   GO OOP: INTERFACES & INTERFACE TYPES           ")
	fmt.Println("##################################################################")

	demo1ImplicitImplementation()
	demo2PointerVsValueReceivers()
	demo3InterfaceInternalRepresentation()
	demo4TheNilInterfaceTrap()
	demo5TheEmptyInterface()
	demo6TypeAssertionsAndCapabilityQuery()
	demo7TypeSwitches()
	demo8InterfaceComposition()
	demo9InterfaceTypesMethodSetsVsTypeSets()
	demo10CompileTimeCheckAndDependencyInjection()

	fmt.Println("\n==================================================================")
	fmt.Println("All interface demonstrations completed successfully!")
	fmt.Println("==================================================================")
}

// ============================================================================
// 1. IMPLICIT IMPLEMENTATION & STRUCTURAL POLYMORPHISM
// ============================================================================
// In Go, interfaces are satisfied IMPLICITLY.
// There is NO 'implements' keyword (unlike Java/C#).
// If a concrete type implements all methods declared in an interface signature,
// the compiler automatically recognizes that the type satisfies the interface.
//
// Benefits:
// - Zero coupling between interface and implementation.
// - Third-party types can satisfy your interfaces retroactively!
// - Polymorphic code can operate on diverse types through a shared contract.

type Shape interface {
	Name() string
	Area() float64
	Perimeter() float64
}

// Rectangle satisfies Shape implicitly
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Name() string      { return "Rectangle" }
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Circle satisfies Shape implicitly
type Circle struct {
	Radius float64
}

func (c Circle) Name() string      { return "Circle" }
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Triangle satisfies Shape implicitly
type Triangle struct {
	A, B, C float64 // sides
}

func (t Triangle) Name() string { return "Triangle" }
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}
func (t Triangle) Area() float64 {
	// Heron's formula
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Polymorphic function: consumes the Shape interface
func PrintShapeSummary(s Shape) {
	fmt.Printf("• %-10s | Area: %8.2f | Perimeter: %8.2f\n",
		s.Name(), s.Area(), s.Perimeter())
}

func demo1ImplicitImplementation() {
	fmt.Println("\n==================================================================")
	fmt.Println("1. IMPLICIT IMPLEMENTATION & POLYMORPHISM")
	fmt.Println("==================================================================")
	fmt.Println("Types implement interfaces automatically by implementing their methods.")
	fmt.Println("No 'implements' keyword needed.")
	fmt.Println()

	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 4},
		Triangle{A: 3, B: 4, C: 5},
	}

	for _, s := range shapes {
		PrintShapeSummary(s)
	}
}

// ============================================================================
// 2. POINTER VS. VALUE RECEIVERS & INTERFACE METHOD SETS
// ============================================================================
// Go's type system has strict rules regarding Method Sets:
// - For a value of type T:
//     Method set contains ONLY methods with receiver (t T).
// - For a pointer of type *T:
//     Method set contains methods with receiver (*t T) AND (t T).
//
// CONSEQUENCE:
// If a method has a pointer receiver (*T), then only *T satisfies the interface.
// A plain value of type T DOES NOT satisfy the interface!
//
// WHY?
// When a value is stored in an interface, Go stores a COPY.
// Interface values are NOT addressable. If Go allowed calling a pointer-receiver
// method on a copied value inside an interface, mutations would either modify
// the temporary copy or cause memory errors. Go strictly disallows this.

type Incrementer interface {
	Increment()
	Value() int
}

type Counter struct {
	count int
}

// Pointer receiver: mutates internal state
func (c *Counter) Increment() {
	c.count++
}

// Value receiver: reads internal state
func (c Counter) Value() int {
	return c.count
}

func executeIncrement(inc Incrementer) {
	inc.Increment()
	fmt.Printf("Counter value after increment: %d\n", inc.Value())
}

func demo2PointerVsValueReceivers() {
	fmt.Println("\n==================================================================")
	fmt.Println("2. POINTER VS. VALUE RECEIVERS & INTERFACE METHOD SETS")
	fmt.Println("==================================================================")

	// Valid: *Counter has both (c *Counter) and (c Counter) methods in its method set.
	cPtr := &Counter{count: 10}
	var inc Incrementer = cPtr
	executeIncrement(inc)

	// What happens with value Counter?
	// cVal := Counter{count: 10}
	// var inc2 Incrementer = cVal
	//
	// COMPILE ERROR:
	// cannot use cVal (variable of type Counter) as Incrementer value in variable declaration:
	// Counter does not implement Incrementer (method Increment has pointer receiver)
	fmt.Println("\nMethod Set Rule:")
	fmt.Println("• *Counter implements Incrementer  [VALID: pointer has both *T and T methods]")
	fmt.Println("• Counter does NOT implement Incrementer [INVALID: method Increment has pointer receiver]")
	fmt.Println("-> Values inside interfaces are not addressable, preventing pointer methods on value copies.")
}

// ============================================================================
// 3. INTERFACE INTERNAL REPRESENTATION (iface & eface)
// ============================================================================
// In the Go runtime, an interface value is a TWO-WORD "FAT POINTER":
//
// 1. Non-empty interface (iface):
//    [ tab *itab ] ----> itab contains:
//                         - Interface type metadata
//                         - Concrete type metadata (hash, size, type descriptor)
//                         - Function pointers to concrete methods
//    [ data unsafe.Pointer ] ----> points to the concrete value in memory
//
// 2. Empty interface (eface, any):
//    [ _type *_type ] ----> Concrete type metadata
//    [ data  unsafe.Pointer ] ----> points to concrete value
//
// An interface is only 16 bytes on 64-bit architectures (2 machine words).

func demo3InterfaceInternalRepresentation() {
	fmt.Println("\n==================================================================")
	fmt.Println("3. INTERFACE INTERNAL REPRESENTATION (iface: tab + data)")
	fmt.Println("==================================================================")

	var s Shape // Zero value of interface

	fmt.Println("--- Zero Value of an Interface ---")
	fmt.Printf("Value: %v | Dynamic Type: %T | Is nil? %t\n", s, s, s == nil)
	fmt.Println("A zero interface has (itab = nil, data = nil).")

	fmt.Println("\n--- Assigning Concrete Types to Interface ---")
	s = Circle{Radius: 7.5}
	fmt.Printf("Assigned Circle -> Dynamic Type: %-15T | Dynamic Value: %v\n", s, s)
	fmt.Printf("Reflection inspection -> Type: %s, Kind: %s\n",
		reflect.TypeOf(s), reflect.ValueOf(s).Kind())

	s = &Rectangle{Width: 20, Height: 10}
	fmt.Printf("Assigned *Rect  -> Dynamic Type: %-15T | Dynamic Value: %v\n", s, s)
	fmt.Printf("Reflection inspection -> Type: %s, Kind: %s\n",
		reflect.TypeOf(s), reflect.ValueOf(s).Kind())

	// Interface variable memory footprint:
	var emptyShape Shape
	fmt.Printf("\nInterface variable memory footprint: %d bytes (2 machine words: tab + data)\n",
		unsafe.Sizeof(emptyShape))
}

// ============================================================================
// 4. THE INFAMOUS "NIL INTERFACE" GOTCHA (TYPED NIL POINTER)
// ============================================================================
// This is one of the most infamous bugs in Go!
//
// RULE: An interface evaluates to nil (i == nil) IF AND ONLY IF:
//       both its dynamic type is nil AND its dynamic value is nil.
//
// If you assign a typed nil pointer to an interface:
//       dynamic type = *CustomError (NOT NIL!)
//       dynamic value = nil
//
// Result: (err != nil) is TRUE! Even though the underlying pointer is nil!

type CustomAppError struct {
	Code    int
	Message string
}

func (e *CustomAppError) Error() string {
	if e == nil {
		return "<nil CustomAppError>"
	}
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// BUGGY PATTERN: returns concrete pointer typed variable directly as error interface.
func findUserBad(id string) error {
	var err *CustomAppError = nil // typed pointer initialized to nil
	if id == "" {
		err = &CustomAppError{Code: 400, Message: "empty ID"}
		return err
	}
	// BUG: Returning 'err' (*CustomAppError) wraps a nil pointer in an error interface!
	// Interface becomes: (type = *CustomAppError, value = nil)
	// That interface itself is NOT nil!
	return err
}

// IDIOMATIC PATTERN: return nil literal directly when there is no error.
func findUserGood(id string) error {
	if id == "" {
		return &CustomAppError{Code: 400, Message: "empty ID"}
	}
	// Explicitly return untyped 'nil'. Interface becomes: (type = nil, value = nil).
	return nil
}

func demo4TheNilInterfaceTrap() {
	fmt.Println("\n==================================================================")
	fmt.Println("4. THE INFAMOUS NIL INTERFACE TRAP (Typed Nil Pointer vs. Nil Interface)")
	fmt.Println("==================================================================")

	fmt.Println("Testing BUGGY implementation (findUserBad):")
	errBad := findUserBad("valid-user-123")
	if errBad != nil {
		fmt.Printf("❌ TRAP TRIGGERED! errBad != nil is TRUE! (Dynamic Type: %T, Dynamic Value: %v)\n",
			errBad, errBad)
		fmt.Println("   The caller thinks an error occurred, even though the operation succeeded!")
	} else {
		fmt.Println("   Success: errBad is nil")
	}

	fmt.Println("\nTesting IDIOMATIC implementation (findUserGood):")
	errGood := findUserGood("valid-user-123")
	if errGood != nil {
		fmt.Println("   Error occurred")
	} else {
		fmt.Printf("✅ CORRECT: errGood == nil is TRUE! (Dynamic Type: %T, Dynamic Value: %v)\n",
			errGood, errGood)
	}

	fmt.Println("\n-> GOLDEN RULE: When returning an interface, always return the explicit literal 'nil',")
	fmt.Println("   never return a typed pointer variable holding nil!")
}

// ============================================================================
// 5. THE EMPTY INTERFACE (any / interface{})
// ============================================================================
// 'interface{}' has zero methods, which means ALL types satisfy it.
// In Go 1.18+, 'any' was added as a predeclared alias for 'interface{}'.
//
// Key nuances:
// 1. Used for dynamic/heterogeneous collections or universal functions (fmt.Println).
// 2. []int CANNOT be passed as []any! Why?
//    A []int is a slice header pointing to contiguous 8-byte ints.
//    A []any is a slice header pointing to contiguous 16-byte (iface) fat pointers.
//    Their memory representations are completely different.

func printAnySlice(name string, items []any) {
	fmt.Printf("%s ([]any of length %d): [ ", name, len(items))
	for _, item := range items {
		fmt.Printf("(%T: %v) ", item, item)
	}
	fmt.Println("]")
}

func demo5TheEmptyInterface() {
	fmt.Println("\n==================================================================")
	fmt.Println("5. THE EMPTY INTERFACE (any / interface{})")
	fmt.Println("==================================================================")

	// Heterogeneous slice using 'any'
	mixed := []any{
		"Hello, Go!",
		42,
		3.14159,
		true,
		Rectangle{Width: 4, Height: 3},
	}
	printAnySlice("Heterogeneous items", mixed)

	// Converting a concrete slice []int to []any requires explicit element copying:
	numbers := []int{10, 20, 30, 40}
	// printAnySlice("Numbers", numbers) // COMPILE ERROR: cannot use numbers (type []int) as []any
	anyNumbers := make([]any, len(numbers))
	for i, v := range numbers {
		anyNumbers[i] = v
	}
	printAnySlice("Converted []int -> []any", anyNumbers)
	fmt.Println("-> []T cannot be passed as []any directly due to differing memory layouts.")
}

// ============================================================================
// 6. TYPE ASSERTIONS & CAPABILITY QUERYING
// ============================================================================
// A Type Assertion retrieves the concrete value from an interface:
// 1. Unsafe assertion: v := i.(ConcreteType) -> PANICS if types do not match!
// 2. Safe "comma-ok" assertion: v, ok := i.(ConcreteType) -> returns ok = false on mismatch.
// 3. Capability Querying (Interface-to-Interface):
//    Tests if an underlying value in interface A also implements interface B!

type Flusher interface {
	Flush() string
}

type BufferWithFlush struct {
	bytes.Buffer
}

func (b *BufferWithFlush) Flush() string {
	content := b.String()
	b.Reset()
	return fmt.Sprintf("Flushed %d bytes: '%s'", len(content), content)
}

func demo6TypeAssertionsAndCapabilityQuery() {
	fmt.Println("\n==================================================================")
	fmt.Println("6. TYPE ASSERTIONS & CAPABILITY QUERYING")
	fmt.Println("==================================================================")

	var val any = "Learning Go Interfaces"

	// 1. Safe comma-ok type assertion
	if str, ok := val.(string); ok {
		fmt.Printf("Safe assertion succeeded: string value = %q (len=%d)\n", str, len(str))
	}

	// 2. Safe assertion failure
	if num, ok := val.(int); !ok {
		fmt.Println("Safe assertion handled failure cleanly: val is NOT an int")
		_ = num
	}

	// 3. Unsafe assertion with panic recovery demonstration
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Unsafe assertion caught panic as expected: %v\n", r)
			}
		}()
		_ = val.(float64) // Will panic!
	}()

	// 4. Capability Querying (Interface-to-Interface):
	// A standard io.Writer can be dynamically checked for additional interfaces (e.g. Flusher)
	fmt.Println("\n--- Capability Querying (Interface to Interface) ---")
	var w io.Writer = &BufferWithFlush{}
	w.Write([]byte("Important buffered message"))

	// Query: Does this io.Writer also implement Flusher?
	if flusher, ok := w.(Flusher); ok {
		fmt.Println("Capability detected: Writer implements Flusher!")
		fmt.Printf("Result: %s\n", flusher.Flush())
	} else {
		fmt.Println("Capability not supported")
	}
}

// ============================================================================
// 7. TYPE SWITCHES
// ============================================================================
// A Type Switch compares an interface's dynamic type against multiple types.
// Inside each case block, the switch variable is automatically typed to the matched type!

func inspectDynamicType(x any) {
	switch v := x.(type) {
	case nil:
		fmt.Println("Type: nil (no dynamic type or value)")
	case int:
		fmt.Printf("Type: int | Squared: %d\n", v*v)
	case string:
		fmt.Printf("Type: string | Uppercase: %s\n", strings.ToUpper(v))
	case bool:
		fmt.Printf("Type: bool | Negated: %t\n", !v)
	case Shape:
		fmt.Printf("Type: Shape (%s) | Area: %.2f\n", v.Name(), v.Area())
	case []int:
		fmt.Printf("Type: []int | Length: %d, Sum: %d\n", len(v), sumInts(v))
	default:
		fmt.Printf("Type: unknown concrete type (%T) | Value: %v\n", v, v)
	}
}

func sumInts(s []int) int {
	sum := 0
	for _, n := range s {
		sum += n
	}
	return sum
}

func demo7TypeSwitches() {
	fmt.Println("\n==================================================================")
	fmt.Println("7. TYPE SWITCHES")
	fmt.Println("==================================================================")

	inspectDynamicType(nil)
	inspectDynamicType(7)
	inspectDynamicType("concurrency")
	inspectDynamicType(true)
	inspectDynamicType(Rectangle{Width: 6, Height: 4})
	inspectDynamicType([]int{1, 2, 3, 4, 5})
	inspectDynamicType(complex(2, 3))
}

// ============================================================================
// 8. INTERFACE COMPOSITION & EMBEDDING
// ============================================================================
// Go encourages small, single-method interfaces:
// - "The bigger the interface, the weaker the abstraction." (Rob Pike)
// Larger interfaces are built by COMPOSING smaller ones through interface embedding.
//
// Examples from the standard library:
// - io.Reader:  Read(p []byte) (n int, err error)
// - io.Writer:  Write(p []byte) (n int, err error)
// - io.Closer:  Close() error
// - io.ReadWriter: embeds Reader + Writer
// - io.ReadWriteCloser: embeds Reader + Writer + Closer

type SimpleReader interface {
	ReadData() string
}

type SimpleWriter interface {
	WriteData(data string)
}

type SimpleCloser interface {
	CloseConnection() error
}

// Composed Interface: embeds SimpleReader and SimpleWriter
type SimpleReadWriter interface {
	SimpleReader
	SimpleWriter
}

// Composed Interface: embeds SimpleReadWriter and SimpleCloser
type SimpleReadWriteCloser interface {
	SimpleReadWriter
	SimpleCloser
}

// InMemoryPipe satisfies SimpleReadWriteCloser implicitly
type InMemoryPipe struct {
	buffer string
	closed bool
}

func (p *InMemoryPipe) ReadData() string {
	if p.closed {
		return "<closed>"
	}
	out := p.buffer
	p.buffer = ""
	return out
}

func (p *InMemoryPipe) WriteData(data string) {
	if !p.closed {
		p.buffer += data
	}
}

func (p *InMemoryPipe) CloseConnection() error {
	p.closed = true
	return nil
}

func demo8InterfaceComposition() {
	fmt.Println("\n==================================================================")
	fmt.Println("8. INTERFACE COMPOSITION / EMBEDDING")
	fmt.Println("==================================================================")

	pipe := &InMemoryPipe{}

	// Can be used as SimpleWriter
	var writer SimpleWriter = pipe
	writer.WriteData("Streamed message through composed interface")

	// Can be used as SimpleReader
	var reader SimpleReader = pipe
	fmt.Printf("Read from pipe: %q\n", reader.ReadData())

	// Can be used as fully composed SimpleReadWriteCloser
	var rwc SimpleReadWriteCloser = pipe
	rwc.WriteData("Closing soon...")
	fmt.Printf("Read before close: %q\n", rwc.ReadData())
	_ = rwc.CloseConnection()
	fmt.Printf("Read after close: %q\n", rwc.ReadData())
	fmt.Println("-> Composed interfaces allow combining focused capabilities seamlessly.")
}

// ============================================================================
// 9. MODERN GO INTERFACE TYPES: METHOD SETS VS TYPE SETS (Go 1.18+ Generics)
// ============================================================================
// Since Go 1.18, interfaces define TYPE SETS, not just method sets:
//
// 1. BASIC INTERFACES:
//    - Contain ONLY methods.
//    - Can be used as regular variable types, parameter types, struct fields.
//
// 2. GENERAL INTERFACES (Type Sets with Type Elements):
//    - Contain type elements (e.g. 'int | float64', '~string', 'comparable').
//    - Can ONLY be used as type constraints in generic type parameters!
//    - CANNOT be used as normal variable types.

// General Interface (Constraint Interface)
// The '~' operator indicates approximation: matches any type whose UNDERLYING type is int/float64.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Generic function constrained by Number interface
func Sum[T Number](values ...T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// Custom defined types with underlying float64/int
type Celsius float64
type Score int

func demo9InterfaceTypesMethodSetsVsTypeSets() {
	fmt.Println("\n==================================================================")
	fmt.Println("9. INTERFACE TYPES: METHOD SETS VS TYPE SETS (Generics Constraints)")
	fmt.Println("==================================================================")

	// Using the generic function with primitive types
	intSum := Sum(10, 20, 30, 40)
	floatSum := Sum(1.5, 2.5, 3.5)
	fmt.Printf("Sum of ints:    %d\n", intSum)
	fmt.Printf("Sum of floats:  %.2f\n", floatSum)

	// Approximation operator (~) allows custom types with matching underlying types:
	celsiusTemps := []Celsius{23.5, 25.0, 21.2}
	celsiusSum := Sum(celsiusTemps...)
	fmt.Printf("Sum of Celsius: %.2f°C (underlying type ~float64 satisfied!)\n", celsiusSum)

	scoreTotal := Sum(Score(100), Score(250), Score(175))
	fmt.Printf("Sum of Scores:  %d (underlying type ~int satisfied!)\n", scoreTotal)

	// Note on compile restriction:
	// var n Number = 42
	// COMPILE ERROR: cannot use type Number outside a type constraint: interface contains type constraints
	fmt.Println("\nRule for General Interfaces:")
	fmt.Println("• Basic Interfaces (only methods) -> usable anywhere (variables, params, fields)")
	fmt.Println("• General Interfaces (with type sets/unions) -> ONLY usable as Generic type constraints!")
}

// ============================================================================
// 10. COMPILE-TIME VERIFICATION & CLIENT-DEFINED INTERFACES (MOCKING)
// ============================================================================
// Best Practice 1: Compile-time check
//   var _ UserStore = (*MemoryStore)(nil)
//
// Best Practice 2: Client-Defined Interfaces & Dependency Injection
//   Package 'service' defines what it needs (UserStore).
//   We can effortlessly plug in MemoryStore (prod/dev) or MockStore (testing).

func demo10CompileTimeCheckAndDependencyInjection() {
	fmt.Println("\n==================================================================")
	fmt.Println("10. IDIOMATIC PATTERNS: COMPILE-TIME CHECKS & DEPENDENCY INJECTION")
	fmt.Println("==================================================================")

	// 1. Production / Dev execution using MemoryStore
	fmt.Println("--- Running UserService with MemoryStore ---")
	memStore := service.NewMemoryStore()
	userSvc := service.NewUserService(memStore)

	user1, err := userSvc.Register("usr-101", "gopher_master", "gopher@golang.org")
	if err != nil {
		fmt.Printf("Registration error: %v\n", err)
		return
	}
	fmt.Printf("Registered user: ID=%s, Username=%s, Email=%s\n",
		user1.ID, user1.Username, user1.Email)

	// Test duplicate email check
	_, err = userSvc.Register("usr-102", "imposter", "gopher@golang.org")
	fmt.Printf("Duplicate registration prevented: %v\n", err)

	// 2. Unit Testing / Mock execution using MockStore
	fmt.Println("\n--- Testing UserService with MockStore (No Database Required!) ---")
	mockStore := service.NewMockStore()

	// Customize mock behavior for FindByEmail
	mockStore.FindByEmailFn = func(email string) (*service.User, error) {
		// Simulate that no user exists with this email
		return nil, service.ErrUserNotFound
	}

	mockSvc := service.NewUserService(mockStore)
	mockUser, err := mockSvc.Register("usr-999", "tester", "test@domain.com")
	if err != nil {
		fmt.Printf("Mock registration error: %v\n", err)
		return
	}

	fmt.Printf("Mock user registered: %s\n", mockUser.Username)
	fmt.Printf("Mock recorded saved users count: %d\n", len(mockStore.SavedUsers))
	fmt.Println("-> Decoupled interface allows 100% test coverage without running external databases.")
}
