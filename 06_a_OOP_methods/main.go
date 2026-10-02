package main

import (
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// 1. METHOD DECLARATIONS
// ============================================================================
// In Go, an "object" is simply a value or variable of some type, and a "method"
// is a function associated with a particular type.
//
// Syntax:
//   func (recv ReceiverType) MethodName(parameters) (returns)
//
// The receiver parameter is placed between the 'func' keyword and the method name.
// Go convention: use a short, 1-2 letter name for the receiver (e.g., 'p' for Point),
// NOT generic names like 'this' or 'self'.

// Point represents a 2D coordinate in space.
type Point struct {
	X, Y float64
}

// Distance is a METHOD with a value receiver of type Point.
// Notice that 'p' is passed by value (a copy is created).
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Regular standalone function for comparison:
// Notice: DistanceFunc(p, q) vs p.Distance(q)
func DistanceFunc(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Rectangle represents a 2D box defined by two corner points.
type Rectangle struct {
	Min, Max Point
}

// Area calculates the surface area of the rectangle.
func (r Rectangle) Area() float64 {
	width := r.Max.X - r.Min.X
	height := r.Max.Y - r.Min.Y
	return math.Abs(width * height)
}

// Perimeter calculates the perimeter length.
func (r Rectangle) Perimeter() float64 {
	width := math.Abs(r.Max.X - r.Min.X)
	height := math.Abs(r.Max.Y - r.Min.Y)
	return 2 * (width + height)
}

// ----------------------------------------------------------------------------
// Methods on Non-Struct Defined Types
// ----------------------------------------------------------------------------
// Methods can be declared on ANY named type created in the same package,
// not just structs! (Except pointer types or interface types).

// Celsius represents temperature in degrees Celsius.
type Celsius float64

// Fahrenheit represents temperature in degrees Fahrenheit.
type Fahrenheit float64

// ToFahrenheit converts a Celsius temperature to Fahrenheit.
// Celsius is a primitive float64 under the hood, but it can have methods!
func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

// String implements the fmt.Stringer interface for pretty printing.
func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", c)
}

// Path is a slice of Points representing a route or polyline.
type Path []Point

// Distance computes the total length along the path.
// Receiver is a defined slice type!
func (path Path) Distance() float64 {
	var total float64
	for i := range path {
		if i > 0 {
			total += path[i-1].Distance(path[i])
		}
	}
	return total
}

// ============================================================================
// 2. METHODS WITH POINTER RECEIVERS
// ============================================================================
// When a method needs to mutate the receiver's state, or when the receiver
// is a large struct that would be expensive to copy, use a POINTER RECEIVER (*T).
//
// Rule of thumb:
// 1. If any method of a type requires a pointer receiver, make ALL methods
//    pointer receivers for consistency.
// 2. If the method mutates the receiver, it MUST be a pointer receiver.
// 3. If the type contains synchronization primitives (e.g. sync.Mutex),
//    it MUST use pointer receivers to avoid copying the lock.

// Offset shifts the point by (dx, dy) using a VALUE receiver.
// BUG/LIMITATION: Because 'p' is a copy, caller's original point will NOT change!
func (p Point) Offset(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// Scale scales the point coordinates by a factor using a POINTER receiver.
// This correctly mutates the original struct caller.
func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

// Translate moves the point by (dx, dy) using a POINTER receiver.
func (p *Point) Translate(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// LargeData simulates a heavy payload where passing by value is wasteful.
type LargeData struct {
	ID      int
	Payload [1024]byte // 1 KB payload
}

// ProcessByPointer avoids copying the 1024-byte payload.
func (ld *LargeData) ProcessByPointer() int {
	return len(ld.Payload)
}

// ============================================================================
// 3. NIL AS A RECEIVER VALUE
// ============================================================================
// In Go, calling a method on a nil pointer receiver DOES NOT automatically panic!
// In languages like Java or C++, `ptr.method()` on null immediately throws
// NullPointerException before method execution.
// In Go, the method is dispatched normally and the receiver variable receives 'nil'.
//
// The panic only occurs if the method body tries to dereference the nil pointer
// (e.g., accessing fields without checking for nil).
//
// Idiomatic Go leverages this to design elegant, nil-safe data structures.

// IntList is a singly-linked list of integers.
// An empty list is canonically represented by a nil pointer (*IntList(nil))!
type IntList struct {
	Value int
	Tail  *IntList
}

// Sum returns the sum of elements in the list.
// NIL RECEIVER IS COMPLETELY VALID: an empty list has a sum of 0!
func (list *IntList) Sum() int {
	if list == nil {
		return 0
	}
	return list.Value + list.Tail.Sum()
}

// Len returns the count of items in the list.
// NIL RECEIVER IS VALID: an empty list has length 0!
func (list *IntList) Len() int {
	if list == nil {
		return 0
	}
	return 1 + list.Tail.Len()
}

// String provides a clean string representation.
func (list *IntList) String() string {
	if list == nil {
		return "nil"
	}
	var b strings.Builder
	b.WriteString("[")
	curr := list
	for curr != nil {
		b.WriteString(fmt.Sprintf("%d", curr.Value))
		if curr.Tail != nil {
			b.WriteString(" -> ")
		}
		curr = curr.Tail
	}
	b.WriteString("]")
	return b.String()
}

// ServerConfig demonstrates a nil-safe configuration getter pattern.
type ServerConfig struct {
	Host string
	Port int
}

// GetHost returns the configured host or a safe default if config is nil.
func (c *ServerConfig) GetHost() string {
	if c == nil || c.Host == "" {
		return "127.0.0.1" // Fallback default
	}
	return c.Host
}

// GetPort returns the configured port or a safe default if config is nil.
func (c *ServerConfig) GetPort() int {
	if c == nil || c.Port == 0 {
		return 8080 // Fallback default
	}
	return c.Port
}

// UnsafeNode demonstrates an unguarded method that panics when receiver is nil.
type UnsafeNode struct {
	Data string
}

func (u *UnsafeNode) BadAccess() string {
	// Dereferencing without checking `u == nil` will cause panic!
	return u.Data
}

// ============================================================================
// DEMO RUNNERS
// ============================================================================

func demoMethodDeclarations() {
	fmt.Println("==================================================================")
	fmt.Println("1. METHOD DECLARATIONS & RECEIVER BASICS")
	fmt.Println("==================================================================")

	p := Point{X: 1, Y: 2}
	q := Point{X: 4, Y: 6}

	// Calling method vs calling function
	distMethod := p.Distance(q)
	distFunc := DistanceFunc(p, q)
	fmt.Printf("Point p: %+v, Point q: %+v\n", p, q)
	fmt.Printf("Method call (p.Distance(q)):   %.4f\n", distMethod)
	fmt.Printf("Function call (Distance(p, q)): %.4f\n", distFunc)

	// Methods on another struct: Rectangle
	rect := Rectangle{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 5},
	}
	fmt.Printf("\nRectangle: Min=%+v, Max=%+v\n", rect.Min, rect.Max)
	fmt.Printf("Rectangle Area:      %.2f\n", rect.Area())
	fmt.Printf("Rectangle Perimeter: %.2f\n", rect.Perimeter())

	// Methods on non-struct defined types (Celsius)
	boilingWater := Celsius(100.0)
	fmt.Printf("\nCustom defined type (Celsius): %s\n", boilingWater)
	fmt.Printf("Converted to Fahrenheit:       %.1f°F\n", boilingWater.ToFahrenheit())

	// Methods on defined slice type (Path)
	perimeterPath := Path{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
		{X: 0, Y: 0}, // closed loop triangle (3-4-5)
	}
	fmt.Printf("\nPath route: %v\n", perimeterPath)
	fmt.Printf("Total path distance: %.2f (expected 12.00)\n", perimeterPath.Distance())
}

func demoPointerReceivers() {
	fmt.Println("\n==================================================================")
	fmt.Println("2. METHODS WITH POINTER RECEIVERS & ADDRESSABILITY")
	fmt.Println("==================================================================")

	// Value Receiver vs Pointer Receiver Mutation Test
	pt := Point{X: 3, Y: 4}
	fmt.Printf("Initial point:           %+v\n", pt)

	// Attempting mutation with Value Receiver
	pt.Offset(10, 20)
	fmt.Printf("After pt.Offset(10, 20): %+v (UNCHANGED! Value receiver works on a copy)\n", pt)

	// Mutation with Pointer Receiver
	pt.Translate(10, 20)
	fmt.Printf("After pt.Translate(10,20): %+v (MUTATED! Pointer receiver modified original)\n", pt)

	pt.Scale(2)
	fmt.Printf("After pt.Scale(2):       %+v (MUTATED! Coords doubled)\n", pt)

	// ------------------------------------------------------------------------
	// Syntactic Sugar: Implicit Referencing and Dereferencing
	// ------------------------------------------------------------------------
	fmt.Println("\n--- Go Compiler Syntactic Sugar ---")

	// 1. Calling a Pointer-receiver method on an Addressable Value:
	// 'val' is of type Point, but Scale has receiver '*Point'.
	// Go automatically interprets 'val.Scale(3)' as '(&val).Scale(3)'.
	val := Point{X: 1, Y: 1}
	val.Scale(5) // Equivalent to (&val).Scale(5)
	fmt.Printf("Automatic (&val).Scale(5): %+v\n", val)

	// 2. Calling a Value-receiver method on a Pointer:
	// 'ptr' is of type *Point, but Distance has receiver 'Point'.
	// Go automatically interprets 'ptr.Distance(q)' as '(*ptr).Distance(q)'.
	ptr := &Point{X: 0, Y: 0}
	target := Point{X: 3, Y: 4}
	d := ptr.Distance(target) // Equivalent to (*ptr).Distance(target)
	fmt.Printf("Automatic (*ptr).Distance(target): %.2f\n", d)

	// ------------------------------------------------------------------------
	// Addressability Gotchas
	// ------------------------------------------------------------------------
	fmt.Println("\n--- Addressability Gotchas ---")
	// The automatic '&' only works if the value is ADDRESSABLE in memory!
	// Struct literals and map elements are NOT addressable.
	//
	// The following lines would fail to compile:
	//   Point{1, 2}.Scale(2)       // COMPILE ERROR: cannot call pointer method Scale on Point literal
	//   pointMap["key"].Scale(2)   // COMPILE ERROR: cannot take the address of map element

	// Workaround 1: Assign literal to a variable first
	temp := Point{1, 2}
	temp.Scale(2)
	fmt.Printf("Assigned to variable first: %+v\n", temp)

	// Workaround 2: Use pointer values in maps
	pointMap := map[string]*Point{
		"origin": {X: 0, Y: 0},
		"marker": {X: 5, Y: 5},
	}
	pointMap["marker"].Scale(10) // Valid! Map stores *Point, so it is already a pointer
	fmt.Printf("Pointer stored in map mutated directly: %+v\n", pointMap["marker"])
}

func demoNilReceivers() {
	fmt.Println("\n==================================================================")
	fmt.Println("3. NIL AS A RECEIVER VALUE")
	fmt.Println("==================================================================")

	// ------------------------------------------------------------------------
	// Example A: Singly-Linked List where nil is an empty list
	// ------------------------------------------------------------------------
	fmt.Println("--- Example A: Linked List with nil as Empty List ---")

	var emptyList *IntList // explicitly nil
	fmt.Printf("emptyList variable is nil? %t\n", emptyList == nil)
	fmt.Printf("emptyList.Len():           %d\n", emptyList.Len())
	fmt.Printf("emptyList.Sum():           %d\n", emptyList.Sum())
	fmt.Printf("emptyList.String():        %s\n", emptyList.String())

	// Build a non-empty list: 10 -> 20 -> 30 -> nil
	list := &IntList{
		Value: 10,
		Tail: &IntList{
			Value: 20,
			Tail: &IntList{
				Value: 30,
				Tail:  nil, // tail is nil!
			},
		},
	}
	fmt.Printf("\nPopulated list: %s\n", list.String())
	fmt.Printf("list.Len():     %d\n", list.Len())
	fmt.Printf("list.Sum():     %d\n", list.Sum())

	// ------------------------------------------------------------------------
	// Example B: Fallback defaults for optional configurations
	// ------------------------------------------------------------------------
	fmt.Println("\n--- Example B: Nil-Safe Configuration Getters ---")

	var nilConfig *ServerConfig = nil
	fmt.Printf("nilConfig.GetHost(): %s (safe fallback!)\n", nilConfig.GetHost())
	fmt.Printf("nilConfig.GetPort(): %d (safe fallback!)\n", nilConfig.GetPort())

	customConfig := &ServerConfig{Host: "api.production.internal", Port: 9000}
	fmt.Printf("customConfig.GetHost(): %s\n", customConfig.GetHost())
	fmt.Printf("customConfig.GetPort(): %d\n", customConfig.GetPort())

	// ------------------------------------------------------------------------
	// Example C: What causes a Panic with nil receivers?
	// ------------------------------------------------------------------------
	fmt.Println("\n--- Example C: Guarded vs Unguarded Nil Receivers ---")
	var unsafe *UnsafeNode = nil

	// Safe handling using recover to demonstrate the runtime panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("CAUGHT EXPECTED PANIC on dereference: %v\n", r)
				fmt.Println("-> Calling the method did NOT panic; the PANIC happened INSIDE when accessing 'u.Data'!")
			}
		}()

		fmt.Println("About to call unsafe.BadAccess()...")
		_ = unsafe.BadAccess() // Calls the method without issue, panics inside when reading u.Data
	}()
}

func demoMethodValuesAndExpressions() {
	fmt.Println("\n==================================================================")
	fmt.Println("4. METHOD VALUES VS METHOD EXPRESSIONS")
	fmt.Println("==================================================================")

	p := Point{X: 1, Y: 2}
	q := Point{X: 4, Y: 6}

	// 1. Method Value:
	// Binds a method to a specific receiver instance 'p'.
	// Resulting function takes only the remaining arguments: func(Point) float64.
	distanceFromP := p.Distance
	fmt.Printf("Method Value distanceFromP(q): %.4f\n", distanceFromP(q))

	// 2. Method Expression:
	// Obtains the method from the TYPE itself (Point.Distance).
	// First parameter is explicitly the receiver: func(Point, Point) float64.
	distExp := Point.Distance
	fmt.Printf("Method Expression Point.Distance(p, q): %.4f\n", distExp(p, q))
}

func main() {
	fmt.Println("##################################################################")
	fmt.Println("             GO METHODS & OBJECT-ORIENTED PROGRAMMING             ")
	fmt.Println("##################################################################")

	demoMethodDeclarations()
	demoPointerReceivers()
	demoNilReceivers()
	demoMethodValuesAndExpressions()

	fmt.Println("\nAll method demonstrations completed successfully!")
}
