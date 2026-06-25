# Go Arrays (03_a_arrays)

This project demonstrates the concept and usage of **arrays** in Go.

In Go, an array is a numbered sequence of elements of a single type with a fixed length.

## Key Concepts

### 1. Fixed Size and Type
The length of an array is part of its type. This means `[5]int` and `[10]int` are distinct, incompatible types. You cannot assign or pass a `[5]int` where a `[10]int` is expected.

```go
var a [5]int // a is of type [5]int
```

### 2. Declaration & Zero Values
When declared without explicit initialization, array elements are automatically initialized to their type's **zero value** (e.g., `0` for integers, `""` for strings, `false` for booleans).

```go
var a [5]int          // [0, 0, 0, 0, 0]
var s [3]string       // ["", "", ""]
```

### 3. Initialization Literals
You can initialize an array with values at declaration:
- **Direct definition**: `b := [3]int{10, 20, 30}`
- **Compiler counted**: Using `...` lets the compiler count the elements: `c := [...]int{1, 2, 3}`
- **Sparse definition**: Initialize specific indices: `d := [5]int{1: 12, 4: 42}` (produces `[0, 12, 0, 0, 42]`)

### 4. Value Semantics (Copy Behavior)
Unlike many other languages (like Java, JavaScript, or Python) where arrays are references, **Go arrays are value types**:
- Assigning one array to another creates a **complete copy** of the data.
- Passing an array to a function passes a copy of the array. Modifications inside the function do not affect the caller.
- To modify an array inside a function, you must pass a **pointer** to the array (`*[N]T`).

```go
original := [3]int{1, 2, 3}
copyAr := original // This copies the whole array
copyAr[0] = 999    // original remains [1, 2, 3]
```

### 5. Multi-dimensional Arrays
You can nest arrays to create multi-dimensional structures:
```go
var matrix [2][3]int // A 2-row, 3-column integer matrix
```

### 6. Iterating Over Arrays
You can loop over arrays using traditional `for` loops or the `range` keyword:
```go
fruits := [...]string{"Apple", "Banana"}

// Using range (index and value)
for index, value := range fruits {
    fmt.Println(index, value)
}
```

### 7. Comparison
Arrays of the same type are **comparable** using `==` and `!=`. They are equal if all their corresponding elements are equal.
```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
fmt.Println(a == b) // true
```

---

## Running the Program

Navigate to this directory (`03_a_arrays`) and run:

```bash
go run main.go
```
