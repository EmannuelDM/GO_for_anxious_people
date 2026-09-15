# Working with JSON in Go (04_json)

This project demonstrates how to work with **JSON** (JavaScript Object Notation) in Go using the standard library's [`encoding/json`](https://pkg.go.dev/encoding/json) package.

JSON is the de-facto standard for data exchange in modern web services, REST APIs, and microservices. Go provides high-performance, reflection-based JSON serialization ("marshaling") and deserialization ("unmarshaling") out of the box without requiring external dependencies.

---

## Key Concepts

### 1. The Exported Fields Rule (Crucial!)
* The `encoding/json` package relies on Go's `reflect` package at runtime.
* **Only exported fields** (fields that start with an **uppercase** letter) can be serialized or deserialized.
* Unexported fields (lowercase) are **silently ignored** by the JSON package, regardless of any struct tags applied to them.

```go
type User struct {
    Name     string `json:"name"`     // Exported: will serialize/deserialize
    password string `json:"password"` // Unexported: SILENTLY IGNORED!
}
```

---

### 2. Marshaling and Unmarshaling

- **Marshaling (`json.Marshal`)**: Converts Go data structures into a compact JSON byte slice (`[]byte`).
  ```go
  data, err := json.Marshal(myStruct)
  ```
- **Pretty-Printing (`json.MarshalIndent`)**: Formats JSON with indentation and newlines, ideal for debugging, CLI output, or config files.
  ```go
  prettyData, err := json.MarshalIndent(myStruct, "", "  ")
  ```
- **Unmarshaling (`json.Unmarshal`)**: Parses a JSON byte slice into a target Go variable. **You must pass a pointer** so `Unmarshal` can mutate the target.
  ```go
  var p Product
  err := json.Unmarshal(data, &p)
  ```

---

### 3. Struct Tags & Directives

Struct tags allow you to customize how struct fields map to JSON keys:

| Tag Directive | Example | Behavior |
| :--- | :--- | :--- |
| **Custom Name** | `` `json:"user_id"` `` | Serializes field as `"user_id"` instead of `"UserID"`. |
| **Omit Empty** | `` `json:"bio,omitempty"` `` | Omits the field from output if it has its zero value (`""`, `0`, `false`, `nil`). |
| **Ignore Field** | `` `json:"-"` `` | Always ignores the field during both marshaling and unmarshaling. |
| **String Coercion** | `` `json:"id,string"` `` | Serializes/deserializes numbers or booleans as JSON strings (e.g. `"42"`). |
| **Raw Tag Key** | `` `json:"-,"` `` | Serializes field with literal name `"-"`. |

---

### 4. Handling Dynamic / Arbitrary JSON

When the incoming JSON schema is unknown, variable, or unstructured, decode into a `map[string]any` (or `map[string]interface{}`):

```go
var data map[string]any
if err := json.Unmarshal([]byte(rawJSON), &data); err != nil { ... }
```

> [!WARNING]
> **The `float64` Gotcha**: When unmarshaling JSON into `any`/`interface{}`, the Go JSON decoder parses **all JSON numbers as `float64`** by default (e.g. `42` becomes `float64(42)`). Always use safe type assertions or `decoder.UseNumber()`.

---

### 5. Streaming with `json.Decoder` & `json.Encoder`

When reading from an `io.Reader` (such as `http.Request.Body`, a file, or socket) or writing to an `io.Writer` (such as `http.ResponseWriter` or `os.Stdout`), use `json.NewDecoder` and `json.NewEncoder` instead of reading all bytes into memory.

#### Advantages of Decoder/Encoder:
1. **Memory Efficiency**: Streams large payloads in chunks without buffering entire files in RAM.
2. **`DisallowUnknownFields()`**: Strictly validates that incoming JSON contains only fields defined in the target struct. This catches API typos and unexpected fields:
   ```go
   dec := json.NewDecoder(reader)
   dec.DisallowUnknownFields() // Returns error if unexpected keys appear
   ```
3. **`UseNumber()`**: Decodes numbers as [`json.Number`](https://pkg.go.dev/encoding/json#Number) (a string representation) to prevent precision loss on 64-bit integers.

---

### 6. Polymorphic / Deferred Parsing with `json.RawMessage`

In event-driven architectures or webhooks, messages often share a common outer envelope with a dynamic payload:

```go
type EventEnvelope struct {
    EventType string          `json:"event_type"`
    Timestamp int64           `json:"timestamp"`
    Payload   json.RawMessage `json:"payload"` // Raw JSON bytes delayed for 2nd pass
}
```

By using `json.RawMessage`, you can first unmarshal the envelope, inspect `EventType`, and then unmarshal `Payload` into the appropriate specific struct (`UserSignupEvent`, `PaymentReceivedEvent`, etc.).

---

### 7. Custom Marshaling (`json.Marshaler` & `json.Unmarshaler`)

You can override standard serialization behavior by implementing the `json.Marshaler` and `json.Unmarshaler` interfaces:

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

This is ideal for formatting timestamps as `YYYY-MM-DD`, encoding bitmasks, or serializing domain-specific types.

---

### 8. Pointers, Null, and Optional Fields (PATCH Requests)

To distinguish between:
1. **Field omitted** by the client,
2. **Field explicitly set to `null`**, and
3. **Field set to its zero value** (`0`, `""`, `false`),

use **pointer types** (`*string`, `*int`, `*bool`):

```go
type PatchUserRequest struct {
    DisplayName *string `json:"display_name,omitempty"`
    Score       *int    `json:"score,omitempty"`
}
```

* If `"score"` is omitted or `null` $\rightarrow$ `req.Score == nil`.
* If `"score": 0` $\rightarrow$ `req.Score != nil` and `*req.Score == 0`.

---

## Running the Program

Navigate to this directory (`04_json`) and run:

```bash
go run main.go
```
