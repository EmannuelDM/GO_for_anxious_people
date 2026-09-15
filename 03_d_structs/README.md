# Go Structs (03_d_structs)

This project demonstrates the concept, memory layout, usage, and advanced features of **structs** in Go.

In Go, a **struct** (short for "structure") is a user-defined type that groups together zero or more fields of potentially different types into a single logical unit. Structs are the primary mechanism for creating custom data structures and representing objects in Go.

---

## Key Concepts

### 1. Memory Layout (Contiguous Storage)
* Unlike many object-oriented languages where fields are stored as collections of pointers, Go structs are laid out in a **single contiguous block of memory**.
* This makes structs extremely cache-friendly and highly performant.
* Struct size is determined by the sum of its field sizes, plus any padding added by the compiler to ensure field alignment (aligned to multiples of the word size).

---

### 2. Declaration and Initialization

- **Zero-Value Declaration**: Declaring a struct variable without an initializer initializes all its fields to their respective zero values (`0`, `""`, `false`, `nil`, etc.).
  ```go
  var p Person // Name: "", Age: 0
  ```
- **Struct Literals (Named Fields)**: Standard and recommended approach. Fields can be initialized in any order, and omitted fields are set to their zero values.
  ```go
  p := Person{
      Age:  30,
      Name: "Alice",
  }
  ```
- **Struct Literals (Positional)**: Shorthand syntax where field names are omitted. You **must** specify all fields in the exact order they are declared in the struct definition. Adding a field later breaks this code.
  ```go
  p := Person{"Bob", 25}
  ```

---

### 3. Pointers to Structs and Syntactic Sugar

* In Go, passing a struct to a function copies the entire struct (pass-by-value). To modify the original struct or avoid copying overhead for large structs, pass a **pointer** (`*T`) instead.
* **Automatic Dereferencing**: Go automatically dereferences struct pointers when accessing fields. You can write `ptr.FieldName` instead of the cumbersome `(*ptr).FieldName`.
  ```go
  p := &Person{Name: "Alice", Age: 30}
  p.Age = 31 // Go automatically converts this to (*p).Age = 31 under the hood
  ```
* **Using `new()`**: The built-in `new` function allocates memory for a struct, zero-initializes it, and returns a pointer to it.
  ```go
  p := new(Person) // Type is *Person, fields initialized to zero values
  ```

---

### 4. Anonymous Structs

* Anonymous structs are declared and instantiated inline, without a named type identifier.
* They are useful for one-off operations, such as configuration structures, test table-driven test cases, or temporary API request/response JSON shapes.
  ```go
  config := struct {
      Host string
      Port int
  }{
      Host: "localhost",
      Port: 8080,
  }
  ```

---

### 5. Embedded Structs and Field Promotion

* Go does not support traditional inheritance. Instead, it promotes **composition** via struct embedding.
* When a struct embeds another struct type anonymously (by declaring a field with only a type name, not a variable name), all fields of the embedded struct are **promoted** to the outer struct.
  ```go
  type Address struct {
      City string
  }
  
  type Employee struct {
      ID int
      Address // Embedded field
  }
  
  emp := Employee{ID: 101, Address: Address{City: "New York"}}
  fmt.Println(emp.City) // Accessible directly! (Promoted field)
  ```

---

### 6. Struct Comparability

* Structs are **comparable** using the `==` and `!=` operators **if and only if** all of their fields are comparable.
* Fields of basic types (int, float, string, boolean, pointers) are comparable.
* Fields of non-comparable types (slices, maps, functions) make the entire struct non-comparable.
  ```go
  // Compile-time error if Person contains a slice or map!
  if personA == personB { ... }
  ```

---

### 7. Struct Tags & JSON Serialization

* Struct fields can be annotated with string literals called **struct tags** (placed in backticks).
* Struct tags contain metadata that can be read at runtime via the `reflect` package. They are heavily used by JSON, XML, and database libraries.
  ```go
  type User struct {
      Username string `json:"user_name"`
      Email    string `json:"email"`
      Password string `json:"-"`               // Omit from JSON output
      Website  string `json:"website,omitempty"` // Omit if value is empty/zero
  }
  ```

---

## Running the Program

Navigate to this directory (`03_d_structs`) and run:

```bash
go run main.go
```
