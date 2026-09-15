package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// STRUCT DEFINITIONS FOR DEMONSTRATIONS
// ============================================================================

// 1. Basic Struct
type Product struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	Price float64  `json:"price"`
	Tags  []string `json:"tags"`
}

// 2. Struct Tags Demonstrations
type UserProfile struct {
	// "id,string" forces the integer to be serialized/deserialized as a JSON string: "101"
	ID int `json:"id,string"`

	// Rename key from "Username" to "user_name"
	Username string `json:"user_name"`

	// Standard mapping
	Email string `json:"email"`

	// "-" tells Go to completely ignore this field in serialization and deserialization
	Password string `json:"-"`

	// "omitempty" skips serialization if the field has its zero value ("")
	Bio string `json:"bio,omitempty"`

	// Pointers with "omitempty" skip serialization if nil.
	// If non-nil (even pointing to an empty string ""), it will be serialized.
	Website *string `json:"website,omitempty"`

	// Slices with "omitempty" skip serialization if nil or length == 0
	Roles []string `json:"roles,omitempty"`

	// Standard boolean
	IsActive bool `json:"is_active"`

	// Unexported field (starts with lowercase):
	// The encoding/json package CANNOT access unexported fields via reflection,
	// so it is ALWAYS silently ignored regardless of tags!
	internalNotes string
}

// 3. RawMessage for Polymorphic / Deferred Decoding
type EventEnvelope struct {
	EventType string          `json:"event_type"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"` // Delays JSON parsing of this sub-tree
}

type UserSignupEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type PaymentReceivedEvent struct {
	InvoiceID string  `json:"invoice_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

// 4. Custom JSON Serialization (json.Marshaler & json.Unmarshaler)
// Standard time.Time serializes as RFC3339 ("2026-09-14T11:40:00Z").
// Suppose we need a custom "YYYY-MM-DD" format.
type DateOnly time.Time

const dateFormat = "2006-01-02"

// MarshalJSON implements json.Marshaler interface
func (d DateOnly) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("%q", time.Time(d).Format(dateFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	parsed, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format %q, expected YYYY-MM-DD: %w", s, err)
	}
	*d = DateOnly(parsed)
	return nil
}

func (d DateOnly) String() string {
	return time.Time(d).Format(dateFormat)
}

type Employee struct {
	Name     string   `json:"name"`
	HireDate DateOnly `json:"hire_date"`
}

// 5. Handling Optional & Null Values (Pointers)
type PatchUserRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Score       *int    `json:"score,omitempty"`
}

// ============================================================================
// MAIN DEMO RUNNER
// ============================================================================

func main() {
	fmt.Println("==================================================")
	fmt.Println("             WORKING WITH JSON IN GO              ")
	fmt.Println("==================================================")

	demoBasicMarshalUnmarshal()
	demoStructTags()
	demoDynamicJSON()
	demoStreamingAndDecoderOptions()
	demoRawMessagePolymorphism()
	demoCustomMarshaler()
	demoPointersAndNullability()
}

// ----------------------------------------------------------------------------
// 1. Basic Marshaling and Unmarshaling
// ----------------------------------------------------------------------------
func demoBasicMarshalUnmarshal() {
	fmt.Println("\n--- 1. Basic Marshaling & Unmarshaling ---")

	laptop := Product{
		ID:    1,
		Title: "Mechanical Keyboard",
		Price: 89.99,
		Tags:  []string{"peripherals", "hardware", "gaming"},
	}

	// json.Marshal: Produces compact JSON byte slice ([]byte)
	compactJSON, err := json.Marshal(laptop)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}
	fmt.Println("Compact JSON (json.Marshal):")
	fmt.Println(string(compactJSON))

	// json.MarshalIndent: Produces indented, human-readable JSON
	prettyJSON, err := json.MarshalIndent(laptop, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling indent: %v\n", err)
		return
	}
	fmt.Println("\nPretty JSON (json.MarshalIndent):")
	fmt.Println(string(prettyJSON))

	// json.Unmarshal: Parses JSON []byte back into Go struct
	inputJSON := `{"id":2,"title":"USB-C Dock","price":129.50,"tags":["dock","usb-c"]}`
	var newProduct Product
	if err := json.Unmarshal([]byte(inputJSON), &newProduct); err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}
	fmt.Printf("\nUnmarshaled Struct:\n%+v\n", newProduct)
}

// ----------------------------------------------------------------------------
// 2. Struct Tags Mastery
// ----------------------------------------------------------------------------
func demoStructTags() {
	fmt.Println("\n--- 2. Struct Tags Mastery (omitempty, -, string) ---")

	site := "https://golang.org"
	u1 := UserProfile{
		ID:            42,
		Username:      "gopher_fan",
		Email:         "gopher@example.com",
		Password:      "super_secret_123", // Will be omitted by "-"
		Bio:           "",                 // Will be omitted by "omitempty"
		Website:       &site,              // Non-nil pointer, included
		Roles:         []string{},         // Empty slice, omitted by "omitempty"
		IsActive:      true,
		internalNotes: "Admin verified via audit", // Unexported, completely ignored
	}

	bytesData, _ := json.MarshalIndent(u1, "", "  ")
	fmt.Println("Marshaled UserProfile (Notice ID as string, no password/bio/roles):")
	fmt.Println(string(bytesData))

	// Notice during unmarshaling: ID as string in JSON is converted to int in struct!
	incomingJSON := `{
		"id": "99",
		"user_name": "rob_pike",
		"email": "rob@example.com",
		"password": "incoming_password_ignored",
		"bio": "Co-designer of Go",
		"is_active": true
	}`
	var u2 UserProfile
	_ = json.Unmarshal([]byte(incomingJSON), &u2)
	fmt.Printf("\nUnmarshaled UserProfile:\nID: %d (parsed from string), User: %s, Pass in struct: %q, Bio: %q\n",
		u2.ID, u2.Username, u2.Password, u2.Bio)
}

// ----------------------------------------------------------------------------
// 3. Dynamic / Arbitrary JSON (map[string]any)
// ----------------------------------------------------------------------------
func demoDynamicJSON() {
	fmt.Println("\n--- 3. Dynamic JSON (Arbitrary schemas with map[string]any) ---")

	// When you don't know the schema ahead of time or it varies dynamically
	rawJSON := `{
		"server": "api.prod.local",
		"port": 443,
		"tls_enabled": true,
		"metrics": {
			"p99_latency_ms": 12.45,
			"requests_per_sec": 8500
		},
		"tags": ["prod", "edge"]
	}`

	var dynamicMap map[string]any
	if err := json.Unmarshal([]byte(rawJSON), &dynamicMap); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Parsed into map[string]any:")
	for k, v := range dynamicMap {
		fmt.Printf("  Key: %-12s | Type: %-20T | Value: %v\n", k, v, v)
	}

	// CRITICAL GOTCHA: In Go, all numbers decoded into interface{}/any are float64 by default!
	portVal := dynamicMap["port"]
	fmt.Printf("\nType of 'port' in dynamic JSON is: %T (value: %v)\n", portVal, portVal)

	// Safe type assertion with comma-ok idiom
	if portFloat, ok := portVal.(float64); ok {
		portInt := int(portFloat)
		fmt.Printf("Converted 'port' safely to int: %d\n", portInt)
	}

	// Accessing nested map
	if metrics, ok := dynamicMap["metrics"].(map[string]any); ok {
		fmt.Printf("Nested metrics p99 latency: %v ms\n", metrics["p99_latency_ms"])
	}
}

// ----------------------------------------------------------------------------
// 4. Streaming with json.Decoder & json.Encoder + Advanced Options
// ----------------------------------------------------------------------------
func demoStreamingAndDecoderOptions() {
	fmt.Println("\n--- 4. Streaming with Encoder/Decoder & DisallowUnknownFields ---")

	// When reading from io.Reader (HTTP request body, files, socket) or writing to io.Writer,
	// use json.NewDecoder and json.NewEncoder instead of buffering bytes with Marshal/Unmarshal.

	// Scenario A: DisallowUnknownFields (Strict parsing / Catch typos)
	jsonWithTypo := `{"id": 10, "title": "Desk Lamp", "prissse": 25.0}`
	reader := strings.NewReader(jsonWithTypo)

	dec := json.NewDecoder(reader)
	// By default, Go ignores unknown fields. DisallowUnknownFields enforces strictness!
	dec.DisallowUnknownFields()

	var p Product
	err := dec.Decode(&p)
	if err != nil {
		fmt.Printf("Strict decoder caught error as expected: %v\n", err)
	}

	// Scenario B: UseNumber() to prevent float64 precision loss on large 64-bit integers
	bigNumJSON := `{"large_id": 9223372036854775800}`
	numDec := json.NewDecoder(strings.NewReader(bigNumJSON))
	numDec.UseNumber()

	var numResult map[string]any
	_ = numDec.Decode(&numResult)

	// With UseNumber(), numbers are kept as json.Number (a string subtype) instead of float64
	largeID := numResult["large_id"].(json.Number)
	intVal, _ := largeID.Int64()
	fmt.Printf("Decoded without precision loss using UseNumber(): %d (Type: %T)\n", intVal, largeID)

	// Scenario C: json.NewEncoder writing to a stream (e.g. bytes.Buffer or os.Stdout)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ") // Indentation on stream
	_ = enc.Encode(Product{ID: 5, Title: "Wireless Mouse", Price: 39.99, Tags: []string{"wireless"}})

	fmt.Println("Encoded directly to stream buffer:")
	fmt.Print(buf.String())
}

// ----------------------------------------------------------------------------
// 5. Deferred / Polymorphic Parsing with json.RawMessage
// ----------------------------------------------------------------------------
func demoRawMessagePolymorphism() {
	fmt.Println("\n--- 5. Deferred Parsing with json.RawMessage (Polymorphic Events) ---")

	// In real-world microservices / webhooks, messages often have an envelope
	// with a common header and a dynamic payload dependent on 'event_type'.

	events := []string{
		`{
			"event_type": "user_signup",
			"timestamp": 1773489000,
			"payload": {"user_id": "usr_987", "email": "clark.kent@dailyplanet.com"}
		}`,
		`{
			"event_type": "payment_received",
			"timestamp": 1773489050,
			"payload": {"invoice_id": "inv_456", "amount": 149.99, "currency": "USD"}
		}`,
	}

	for i, rawEvent := range events {
		// First pass: Decode the envelope
		var envelope EventEnvelope
		if err := json.Unmarshal([]byte(rawEvent), &envelope); err != nil {
			fmt.Printf("Envelope decode error: %v\n", err)
			continue
		}

		fmt.Printf("Event %d Type: %q (Payload held in RawMessage: %s)\n",
			i+1, envelope.EventType, string(envelope.Payload))

		// Second pass: Decode payload based on event_type
		switch envelope.EventType {
		case "user_signup":
			var signup UserSignupEvent
			_ = json.Unmarshal(envelope.Payload, &signup)
			fmt.Printf("  -> Handling Signup: UserID=%s, Email=%s\n", signup.UserID, signup.Email)

		case "payment_received":
			var payment PaymentReceivedEvent
			_ = json.Unmarshal(envelope.Payload, &payment)
			fmt.Printf("  -> Handling Payment: Invoice=%s, Amount=%.2f %s\n",
				payment.InvoiceID, payment.Amount, payment.Currency)
		}
	}
}

// ----------------------------------------------------------------------------
// 6. Custom Marshaler and Unmarshaler (json.Marshaler, json.Unmarshaler)
// ----------------------------------------------------------------------------
func demoCustomMarshaler() {
	fmt.Println("\n--- 6. Custom Marshaler & Unmarshaler (YYYY-MM-DD Date) ---")

	parsedTime, _ := time.Parse("2006-01-02", "2024-03-15")
	emp := Employee{
		Name:     "Ada Lovelace",
		HireDate: DateOnly(parsedTime),
	}

	// Marshaling: calls HireDate.MarshalJSON()
	data, _ := json.MarshalIndent(emp, "", "  ")
	fmt.Println("Marshaled Employee with custom DateOnly:")
	fmt.Println(string(data))

	// Unmarshaling: calls HireDate.UnmarshalJSON()
	incoming := `{"name": "Grace Hopper", "hire_date": "2023-11-01"}`
	var newEmp Employee
	if err := json.Unmarshal([]byte(incoming), &newEmp); err != nil {
		fmt.Printf("Custom unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("Unmarshaled Employee: Name=%s, HireDate=%s\n", newEmp.Name, newEmp.HireDate)
}

// ----------------------------------------------------------------------------
// 7. Pointers, Null, and Optional Fields
// ----------------------------------------------------------------------------
func demoPointersAndNullability() {
	fmt.Println("\n--- 7. Pointers, Null, and Optional Fields ---")

	// Pointers allow differentiating 3 states:
	// State 1: Field is omitted -> pointer is nil
	// State 2: Field is explicitly null -> pointer is nil
	// State 3: Field is explicitly provided as zero value (e.g. 0 or "") -> pointer is NOT nil!

	// Case A: score explicitly set to 0 (e.g., reset score in an API PATCH)
	jsonScoreZero := `{"score": 0}`
	var req1 PatchUserRequest
	_ = json.Unmarshal([]byte(jsonScoreZero), &req1)

	if req1.Score != nil {
		fmt.Printf("Case A (Explicit 0): Score was provided and is: %d\n", *req1.Score)
	} else {
		fmt.Println("Case A: Score was omitted")
	}

	// Case B: score omitted completely
	jsonOmitted := `{"display_name": "GoExpert"}`
	var req2 PatchUserRequest
	_ = json.Unmarshal([]byte(jsonOmitted), &req2)

	if req2.Score != nil {
		fmt.Printf("Case B: Score was provided: %d\n", *req2.Score)
	} else {
		fmt.Println("Case B (Omitted): Score is nil (not modified)")
	}

	// Case C: Marshaling pointers
	zero := 0
	name := "Alice"
	patch := PatchUserRequest{
		DisplayName: &name,
		Score:       &zero, // Explicit 0 pointer
	}
	patchJSON, _ := json.Marshal(patch)
	fmt.Printf("Case C: Marshaled patch with pointer to 0: %s\n", string(patchJSON))
}
