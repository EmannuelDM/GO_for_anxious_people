# Go Slices (03_b_slices)

This project demonstrates the concept, internals, and usage of **slices** in Go.

In Go, a **slice** is a dynamically-sized, flexible view into the elements of an array. Slices are much more common than arrays in Go programming due to their flexibility and powerful built-in features.

---

## Key Concepts

### 1. The Slice Header (Internals)
A slice does not store any data itself. Instead, it is a small descriptor (called a **slice header**) that contains:
1. **Pointer**: A reference to the first element of an underlying array that the slice accesses.
2. **Length**: The number of elements currently in the slice. Accessible via `len(s)`.
3. **Capacity**: The maximum number of elements the slice can grow to, starting from the first element in the slice. Accessible via `cap(s)`.

```
Slice Header:
+--------------+----------+--------------+
| Pointer (ptr)| Length   | Capacity     |
|   (address)  | (len)    | (cap)        |
+-------+------+----+-----+-------+------+
        |           |             |
        v           |             |
Underlying Array:   |             |
  [x, x, [10, 20, 30, 40], x, x, x]
         ^           ^
         |           |
         +-----------+ (len = 4)
         |                         |
         +-------------------------+ (cap = 6)
```

---

### 2. Slices vs. Arrays

| Feature | Array | Slice |
| :--- | :--- | :--- |
| **Size** | Fixed at declaration | Dynamic (can grow and shrink) |
| **Type Definition** | Size is part of the type (`[5]int`) | Size is not part of the type (`[]int`) |
| **Memory** | Stores elements directly | Points to an underlying array |
| **Assignment / Passing** | Value semantics (copies all elements) | Reference-like semantics (copies the header, shares memory) |
| **Comparison** | Comparable (`==`, `!=`) | Only comparable to `nil` |

---

### 3. Declaration and Initialization

Slices can be created in several ways:

- **Nil Slice**: Declared without initial value.
  ```go
  var s []int // s == nil, len = 0, cap = 0
  ```
- **Slice Literal**:
  ```go
  s := []int{1, 2, 3} // len = 3, cap = 3
  ```
- **Using `make()`**: Used to pre-allocate size/capacity.
  ```go
  s := make([]int, 5)     // len = 5, cap = 5 (all zero-valued)
  s := make([]int, 3, 10) // len = 3, cap = 10
  ```
- **Slicing an Array**:
  ```go
  arr := [5]int{10, 20, 30, 40, 50}
  s := arr[1:4] // s = [20, 30, 40], len = 3, cap = 4
  ```

---

### 4. Slicing Syntax
You can slice an array or an existing slice using the syntax `s[low:high]`:
* `low`: Inclusive starting index (defaults to `0`).
* `high`: Exclusive ending index (defaults to `len(s)`).
* **Three-index slicing**: `s[low:high:max]` sets the capacity of the resulting slice to `max - low` (cannot exceed the original slice's capacity).

---

### 5. Growing Slices with `append`
The built-in `append()` function adds elements to the end of a slice:
* If the underlying array has enough capacity, `append` updates the array and returns a slice pointing to it.
* If capacity is exceeded, Go **allocates a new, larger array**, copies existing data, appends the new element, and returns a new slice pointing to this new array.

```go
s := []int{1, 2}
s = append(s, 3, 4) // Appends multiple elements
s = append(s, anotherSlice...) // Appends another slice using the variadic operator '...'
```

---

### 6. Copying Slices with `copy`
The built-in `copy()` function copies elements from a source slice to a destination slice:
* It only copies elements up to `min(len(dest), len(src))`.
* **Important**: You must pre-allocate the destination slice's length before copying.
```go
src := []int{1, 2, 3}
dest := make([]int, len(src)) // Pre-allocated
copy(dest, src) // copies elements
```

---

### 7. Slices and Functions
* When a slice is passed to a function, the slice header is copied by value.
* Since the copied header contains a pointer to the same underlying array, **modifications to elements are visible to the caller**.
* However, operations that modify the length or capacity (like `append`) will only affect the local copy of the header. To update the caller, you must either:
  1. Return the new slice and assign it back.
  2. Pass a pointer to the slice (`*[]T`).

---

## Common Pitfalls

### Memory Leaks (Large Underlying Arrays)
Since a slice holds a reference to an underlying array, the entire array will remain in memory as long as at least one slice points to it.
If you have a huge array and create a tiny slice from it, the garbage collector cannot reclaim the huge array.
**Solution**: Copy the small slice's elements to a new slice and discard the original.
```go
hugeSlice := getHugeSlice()
smallSlice := make([]int, 2)
copy(smallSlice, hugeSlice[:2]) // Now hugeSlice can be garbage collected!
```

---

## Running the Program

Navigate to this directory (`03_b_slices`) and run:

```bash
go run main.go
```
