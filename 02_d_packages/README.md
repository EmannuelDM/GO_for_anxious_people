# Go Packages and Imports (02_d_packages)

This project demonstrates the concept, structure, and usage of **packages** and **imports** in Go.

Packages are Go's way of organizing and modularizing code. Every Go program is made up of packages, starting with the `main` package.

## Key Concepts

### 1. Package Structure & Rules
- **One Directory = One Package**: A directory in Go corresponds to a package. All `.go` files in a single directory must declare the same package name at the top (e.g., `package mathutils`).
- **File Sharing**: Code within the same package (even in separate files within the same directory) can access each other's functions, variables, and structs directly, without any imports.
- **Modules**: To import local packages, you must define a module using `go.mod` (created using `go mod init <module-name>`).

### 2. Exported vs. Unexported Identifiers (Visibility)
In Go, visibility is determined solely by the capitalization of the identifier's first letter:
- **Exported (Public)**: Starts with an **uppercase** letter (e.g., `Add`, `Circle`, `Radius`). Exported identifiers can be accessed from outside the package.
- **Unexported (Private)**: Starts with a **lowercase** letter (e.g., `multiply`, `color`). Unexported identifiers are only visible within the package where they are declared.

```go
package mathutils

func Add(a, b int) int      { return a + b } // Exported (visible externally)
func multiply(a, b int) int { return a * b } // Unexported (visible only inside mathutils)
```

For structs:
- The struct name itself must be capitalized to be exported (e.g., `Circle`).
- Struct fields must also be capitalized to be accessible outside the package (e.g., `Radius` is exported, `color` is unexported).

### 3. Importing Packages
To use code from another package, you must import it.
- **Standard Library Packages**: Imported by their name (e.g., `"fmt"`, `"math"`).
- **Local Module Packages**: Imported using the module name defined in `go.mod` followed by the directory path (e.g., `"02_d_packages/mathutils"`).

```go
import (
	"fmt"
	"02_d_packages/mathutils"
)
```

### 4. Package Aliasing
You can rename an imported package by providing an alias before the import path. This is useful for avoiding name collisions or shortening long package names.

```go
import str "02_d_packages/stringutils"

// Now you can call: str.Reverse("hello")
```

### 5. Package Initialization (`init()` Function)
- Every package can contain one or more `init()` functions.
- `init()` functions take no arguments and return no values.
- They are executed automatically when the package is first initialized, **before** the `main()` function starts.
- This is commonly used to perform setup tasks or register components.

---

## Running the Program

Navigate to this directory (`02_d_packages`) and run:

```bash
go run main.go
```
