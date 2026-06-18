package main

// We can import standard library packages and our local custom packages.
// The import path for local packages is: "<module_name>/<package_subdirectory_path>".
import (
	"fmt"
	"math"

	// Importing a local package
	"02_d_packages/mathutils"

	// Importing a local package and giving it an alias
	// This is useful to resolve naming conflicts or shorten long package names.
	str "02_d_packages/stringutils"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("               PACKAGES IN GO                     ")
	fmt.Println("==================================================")

	// 1. Using Standard Library Packages
	// We imported "fmt" and "math" above.
	fmt.Println("\n--- 1. Using Standard Library Packages ---")
	fmt.Printf("Pi from 'math' package: %.5f\n", math.Pi)

	// 2. Using Custom Local Packages
	fmt.Println("\n--- 2. Using Custom Local Package (mathutils) ---")
	sum := mathutils.Add(5, 7)
	diff := mathutils.Subtract(10, 4)
	fmt.Printf("mathutils.Add(5, 7) = %d\n", sum)
	fmt.Printf("mathutils.Subtract(10, 4) = %d\n", diff)

	// 3. Exported vs Unexported (Visibility)
	fmt.Println("\n--- 3. Exported vs Unexported Identifiers ---")
	// Creating a Circle using the constructor function
	myCircle := mathutils.NewCircle(5.0, "red")

	// We can access 'Radius' because it starts with a Capital letter (Exported).
	fmt.Printf("Circle radius: %.2f\n", myCircle.Radius)

	// We can call 'Area()' and 'GetColor()' because they are exported methods.
	fmt.Printf("Circle area: %.2f\n", myCircle.Area())
	fmt.Printf("Circle color (via getter): %s\n", myCircle.GetColor())

	// Accessing unexported members:
	// The following lines would fail to compile if uncommented because:
	// - mathutils.multiply is unexported (starts with lowercase 'm')
	// - myCircle.color is an unexported field (starts with lowercase 'c')
	//
	// _ = mathutils.multiply(2, 3) // compilation error: cannot refer to unexported name mathutils.multiply
	// _ = myCircle.color           // compilation error: myCircle.color undefined (cannot refer to unexported field color)
	fmt.Println("Note: Trying to access unexported identifiers (like 'mathutils.multiply' or 'myCircle.color')")
	fmt.Println("      from 'main' package will cause a compilation error at compile time.")

	// 4. Package Aliasing
	// We aliased "02_d_packages/stringutils" as "str" in the import block.
	fmt.Println("\n--- 4. Package Import Aliasing ---")
	original := "Golang Packages"
	reversed := str.Reverse(original)
	fmt.Printf("Original: %q\n", original)
	fmt.Printf("Reversed (using aliased package 'str'): %q\n", reversed)

	// 5. Package Initialization
	fmt.Println("\n--- 5. Package Initialization (init) ---")
	fmt.Println("Observe the program's output at the very beginning.")
	fmt.Println("Before main() runs, all imported packages are initialized.")
	fmt.Println("Their 'init()' functions are executed in dependency order:")
	fmt.Println("1. Packages with no imports are initialized first.")
	fmt.Println("2. Then packages that import them.")
	fmt.Println("3. Multiple 'init()' functions in a single package are executed in the order they appear.")
}
