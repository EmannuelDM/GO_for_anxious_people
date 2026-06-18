package mathutils

import (
	"fmt"
	"math"
)

func init() {
	fmt.Println("[mathutils/geometry.go] Package mathutils (geometry) initialized!")
}

// pi is unexported. It is only accessible within the mathutils package.
// However, since geometry.go and arithmetic.go are in the same package (mathutils),
// code in this file can access unexported identifiers defined in arithmetic.go (like multiply),
// and arithmetic.go could access pi.
const pi = 3.141592653589793

// Circle represents a geometric circle.
// Exported struct. Note that its fields must also be capitalized to be exported!
type Circle struct {
	Radius float64 // Exported field
	color  string  // Unexported field - cannot be accessed outside the mathutils package
}

// NewCircle is a constructor-like function for Circle.
// It is standard Go practice to return a struct or a pointer to a struct.
func NewCircle(radius float64, color string) Circle {
	return Circle{
		Radius: radius,
		color:  color,
	}
}

// Area calculates and returns the area of the circle.
// Exported method.
func (c Circle) Area() float64 {
	// Accessing the unexported field 'color' is fine here because we are inside the package.
	fmt.Printf("[geometry.go] Calculating area for a %s circle...\n", c.color)
	return pi * math.Pow(c.Radius, 2)
}

// GetColor returns the circle's color.
// Since 'color' is unexported, we can provide an exported getter method if external packages need it.
func (c Circle) GetColor() string {
	return c.color
}

// DoubleArea is a demo function showing that files in the same directory (package)
// can access unexported identifiers in other files of the same package.
// Here we call the unexported 'multiply' function from arithmetic.go.
func DoubleArea(c Circle) float64 {
	area := c.Area()
	// multiply is unexported but defined in arithmetic.go (same package).
	// We convert area to int for demo purposes of multiply, then back to float64.
	doubled := multiply(int(area), 2)
	return float64(doubled)
}
