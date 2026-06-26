# Go Maps (03_c_Maps)

This project demonstrates the concept, internals, and usage of **maps** in Go.

In Go, a **map** is a built-in data type that associates keys with values. Under the hood, a map is implemented as a hash table. Maps are extremely efficient for lookups, insertions, and deletions (average O(1) time complexity).

---

## Key Concepts

### 1. Map Internals (Reference-like Semantics)
A map value is a pointer to a `runtime.hmap` struct header. 
* Because a map is a pointer, assigning a map to another variable or passing it to a function copies the pointer, **not** the underlying data.
* Therefore, modifications to a map inside a function are visible to the caller.

---

### 2. Comparable Key Types
Map keys must be of a type for which the comparison operators `==` and `!=` are fully defined.
* **Allowed as keys**: basic types (integers, floats, strings, booleans, pointers), channels, interfaces, and structs containing only comparable fields.
* **Not allowed as keys**: slices, maps, and functions (since they are not comparable).
* Values can be of any type, including slices, other maps, or functions.

---

### 3. Declaration and Initialization

- **Nil Map**: Declared without an initial value. Writing to a nil map **will panic**!
  ```go
  var m map[string]int // m == nil, len = 0
  // m["key"] = 10     // PANIC: assignment to entry in nil map
  ```
- **Map Literal**:
  ```go
  m := map[string]int{
      "Alice": 25,
      "Bob":   30,
  }
  ```
- **Using `make()`**: Pre-allocates a ready-to-use map.
  ```go
  m := make(map[string]int)     // len = 0, ready to write
  m := make(map[string]int, 50) // Pre-allocates space for 50 items (capacity hint)
  ```

---

### 4. Reading, Writing, and the "Comma OK" Idiom

* **Writing/Updating**:
  ```go
  m["key"] = value
  ```
* **Reading**:
  If a key does not exist, Go returns the **zero value** of the value type.
  ```go
  val := m["missing_key"] // returns 0 if value type is int
  ```
* **Checking for Existence (Comma OK Idiom)**:
  To distinguish between a key that isn't present vs a key that exists with a zero value, use the two-value assignment:
  ```go
  val, ok := m["key"]
  if ok {
      // Key exists, val is its value
  } else {
      // Key does not exist
  }
  ```

---

### 5. Deleting Elements
Use the built-in `delete()` function:
* It takes the map and the key to delete.
* If the key is not in the map, `delete` is a silent **no-op** (safe to call).
```go
delete(m, "key")
```

---

### 6. Iterating Over Maps (Randomized Order)
* You can iterate over maps using `for k, v := range m`.
* **Important**: Go intentionally randomizes map iteration order. You cannot rely on keys being returned in a specific order, even across consecutive runs.
* To iterate in a sorted order, you must collect the keys, sort them manually, and then retrieve the values:
```go
var keys []string
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

---

### 7. Thread Safety & Concurrency
* Go maps are **not safe** for concurrent read and write operations.
* Simultaneous read and write, or multiple writes, will result in a fatal runtime error: `fatal error: concurrent map read and map write`.
* **Solution**: Protect access to the map using a `sync.Mutex` or use `sync.Map` for specific use cases.
```go
type SafeMap struct {
    mu sync.Mutex
    m  map[string]int
}

func (s *SafeMap) Set(k string, v int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[k] = v
}
```

---

## Running the Program

Navigate to this directory (`03_c_Maps`) and run:

```bash
go run main.go
```
