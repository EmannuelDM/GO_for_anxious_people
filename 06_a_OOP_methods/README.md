# OOP Methods in Go (06_a_OOP_methods)

In Go, there are no `class` keywords, inheritance hierarchies, or implicit `this`/`self` pointers. Instead, Go provides a clean, composition-based Object-Oriented model where **methods** are simply functions declared with an explicit **receiver parameter**.

This tutorial covers the three pillars of Go method design:
1. **Method Declarations** (syntax, naming, methods on structs & non-struct defined types, method values vs. expressions).
2. **Methods with Pointer Receivers** (mutability, copying costs, automatic compiler referencing/dereferencing, and addressability rules).
3. **Nil as a Receiver Value** (why `nil` method calls work in Go, safe nil receivers, recursive linked lists, and default-fallback getters).

---

## 1. Method Declarations

### Syntax
A method declaration binds a function to a specific type by specifying a **receiver parameter** between the `func` keyword and the method name:

```go
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```

- **Receiver parameter**: `(p Point)` indicates that `Distance` is a method on the `Point` type.
- **Convention**: Use short, 1- or 2-letter receiver names that match the type (e.g., `p` for `Point`, `c` for `Config`, `r` for `Reader`). Avoid generic names like `this` or `self`.

### Methods vs. Functions
A method is essentially syntactic sugar for a regular function where the receiver is passed as the first argument:

| Calling Syntax | Declaration | Namespace |
| :--- | :--- | :--- |
| `p.Distance(q)` | `func (p Point) Distance(q Point) float64` | Bound to type `Point` |
| `DistanceFunc(p, q)` | `func DistanceFunc(p, q Point) float64` | Package-level function |

Because methods are namespaced to their type, different types can declare methods with identical names without conflicts (e.g., `Rectangle.Area()` and `Circle.Area()`).

### Methods on Non-Struct Defined Types
In Go, methods are **not** restricted to structs! You can declare methods on **any named type** created with `type`, as long as the type is defined in the **same package**:

```go
// 1. Primitive base type
type Celsius float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9.0/5.0 + 32.0)
}

func (c Celsius) String() string {
    return fmt.Sprintf("%.1f°C", c)
}

// 2. Slice base type
type Path []Point

func (path Path) Distance() float64 {
    var total float64
    for i := range path {
        if i > 0 {
            total += path[i-1].Distance(path[i])
        }
    }
    return total
}
```

> [!NOTE]
> **Package Restriction**: You cannot declare methods on built-in types directly (e.g., `func (i int) Double() int`), nor on types defined in another package (like `time.Time`). You must define your own named type alias (`type MyInt int`) first.

### Method Values vs. Method Expressions
Go treats methods as first-class functions in two distinct forms:

```go
p := Point{X: 1, Y: 2}
q := Point{X: 4, Y: 6}

// 1. Method Value: bound to instance 'p'
// Type: func(Point) float64
distanceFromP := p.Distance
dist1 := distanceFromP(q)

// 2. Method Expression: unbound, referenced from the TYPE
// Type: func(Point, Point) float64
distFn := Point.Distance
dist2 := distFn(p, q)
```

---

## 2. Methods with a Pointer Receiver

When declaring a method, you choose between a **value receiver** `(p Point)` and a **pointer receiver** `(p *Point)`.

```go
// Value receiver: operates on a COPY of Point
func (p Point) Offset(dx, dy float64) {
    p.X += dx // Modifies local copy only; caller is unchanged!
}

// Pointer receiver: operates on the ORIGINAL Point
func (p *Point) Translate(dx, dy float64) {
    p.X += dx // Modifies original point in-place
    p.Y += dy
}
```

### Compiler Syntactic Sugar (Implicit `&` and `*`)
To keep code clean and readable, the Go compiler automatically handles referencing and dereferencing:

1. **Automatic Address-Of (`&`)**:
   Calling a pointer-receiver method on an addressable value variable:
   ```go
   val := Point{1, 1}
   val.Scale(5) // Compiler rewrites to: (&val).Scale(5)
   ```

2. **Automatic Dereference (`*`)**:
   Calling a value-receiver method on a pointer:
   ```go
   ptr := &Point{1, 1}
   d := ptr.Distance(q) // Compiler rewrites to: (*ptr).Distance(q)
   ```

### The Addressability Rule (Crucial Gotcha!)
The compiler can only take the address `(&x)` if `x` is **addressable** in memory (i.e., stored in a variable).

```go
// ❌ COMPILE ERROR: Cannot take address of a struct literal
Point{1, 2}.Scale(2)

// ✅ Fix: Assign to a variable first
pt := Point{1, 2}
pt.Scale(2)

// ❌ COMPILE ERROR: Map values are not addressable
pointMap := map[string]Point{"origin": {0, 0}}
pointMap["origin"].Scale(2)

// ✅ Fix: Store pointers in the map instead
pointMapPtr := map[string]*Point{"origin": {0, 0}}
pointMapPtr["origin"].Scale(2)
```

### Value vs. Pointer Receiver Decision Guide

| Criteria | Value Receiver `(T)` | Pointer Receiver `(*T)` |
| :--- | :--- | :--- |
| **Needs to mutate state?** | ❌ No (caller won't see changes) | ✅ **Yes** |
| **Large struct size?** | ❌ Incurs copy overhead | ✅ **Zero copy** (pointer size only) |
| **Contains locks/mutexes?** | ❌ Never copy a `sync.Mutex` | ✅ **Required** |
| **Small / Primitive type?** | ✅ Recommended (`time.Time`, `int`) | ❌ Unnecessary indirection |
| **Consistency Rule** | If most methods are value receivers | If **any** method uses pointer, use pointers for **all** |

---

## 3. Nil as a Receiver Value

In languages like Java, Python, or C++, calling a method on a `null` reference raises an immediate exception (`NullPointerException`) before the method body is ever entered.

**In Go, calling a method on a `nil` pointer receiver is completely valid!**

### Why It Works
A method invocation `recv.Method()` is translated by the compiler to a function call where the receiver is the first argument:

$$\text{recv.Method(arg)} \implies \text{(*Type).Method(recv, arg)}$$

When `recv` is `nil`, the method is dispatched normally and entered. The parameter `recv` simply has the value `nil`. A runtime panic only happens if the method body tries to dereference the nil pointer without checking (`recv.field` or `*recv`).

```go
func (u *UnsafeNode) BadAccess() string {
    return u.Data // Panics if u == nil!
}
```

### Idiomatic Nil-Safe Patterns

#### Pattern A: Recursive Data Structures (Canonical Empty Value)
In a linked list or binary tree, `nil` naturally represents the empty list or empty tree. Methods can handle `nil` gracefully without extra wrapper types:

```go
type IntList struct {
    Value int
    Tail  *IntList
}

// Nil receiver is the empty list!
func (list *IntList) Sum() int {
    if list == nil {
        return 0 // Base case: sum of empty list is 0
    }
    return list.Value + list.Tail.Sum()
}

func (list *IntList) Len() int {
    if list == nil {
        return 0 // Base case: length of empty list is 0
    }
    return 1 + list.Tail.Len()
}

var emptyList *IntList // nil
emptyList.Len()        // returns 0 (No panic!)
emptyList.Sum()        // returns 0 (No panic!)
```

#### Pattern B: Default Fallback Getters
When working with optional configuration structs, nil-safe getters provide sensible defaults without requiring caller null-checks:

```go
type ServerConfig struct {
    Host string
    Port int
}

func (c *ServerConfig) GetHost() string {
    if c == nil || c.Host == "" {
        return "127.0.0.1" // Fallback default
    }
    return c.Host
}

func (c *ServerConfig) GetPort() int {
    if c == nil || c.Port == 0 {
        return 8080 // Fallback default
    }
    return c.Port
}

var cfg *ServerConfig = nil
fmt.Println(cfg.GetHost()) // "127.0.0.1"
fmt.Println(cfg.GetPort()) // 8080
```

---

## 4. Running the Code

Execute the demo module using standard Go CLI:

```bash
cd 06_a_OOP_methods
go run main.go
```

To run static analysis and lint checks:

```bash
go vet ./...
```
