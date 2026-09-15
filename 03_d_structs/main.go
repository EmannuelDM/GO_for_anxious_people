package main

import (
	"encoding/json"
	"fmt"
)

// 1. Defining a Struct
// A struct is a sequence of named elements, called fields, each of which has a name and a type.
type Person struct {
	Name string
	Age  int
}

// 2. Struct for Pointer / Modification demo
type Car struct {
	Brand string
	Model string
	Year  int
}

// 3. Embedded Structs / Composition
// Go does not have traditional OOP inheritance. Instead, it uses composition.
type Address struct {
	City    string
	Country string
}

type Employee struct {
	ID   int
	Name string
	// Anonymous field (embedding): we only specify the type, not the field name.
	// This embeds the fields of Address directly into Employee.
	Address
}

// 4. Struct Tags & JSON Serialization
// Struct tags allow adding metadata to fields. These are used by packages (like encoding/json)
// via reflection to alter how fields are serialized/deserialized.
type User struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
	// "-" tells the JSON encoder to omit this field entirely
	Password string `json:"-"`
	// "omitempty" tells the JSON encoder to skip this field if it has its zero value
	Website string `json:"website,omitempty"`
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("                 STRUCTS IN GO                    ")
	fmt.Println("==================================================")

	// --- 1. Declaration and Zero Values ---
	fmt.Println("\n--- 1. Declaration and Zero Values ---")

	// If we declare a struct without initializing it, its fields are initialized
	// to their respective zero values (e.g., "" for string, 0 for int).
	var p1 Person
	fmt.Printf("p1: %+v (Zero value)\n", p1)
	fmt.Printf("p1.Name: %q, p1.Age: %d\n", p1.Name, p1.Age)

	// --- 2. Struct Literals (Initialization) ---
	fmt.Println("\n--- 2. Struct Literals (Initialization) ---")

	// Named field syntax (Recommended): You don't have to specify all fields,
	// and the order of fields doesn't matter. Unspecified fields get their zero values.
	p2 := Person{
		Age:  30,
		Name: "Alice",
	}
	fmt.Printf("p2 (Named fields): %+v\n", p2)

	// Positional syntax: You must specify all fields in the exact order they are defined.
	// This is generally discouraged because adding a field to the struct will break this code.
	p3 := Person{"Bob", 25}
	fmt.Printf("p3 (Positional fields): %+v\n", p3)

	// --- 3. Accessing and Modifying Fields ---
	fmt.Println("\n--- 3. Accessing and Modifying Fields ---")

	// Fields are accessed and modified using the dot (.) operator.
	var p4 Person
	p4.Name = "Charlie"
	p4.Age = 40
	fmt.Printf("p4 modified: Name=%s, Age=%d\n", p4.Name, p4.Age)

	// --- 4. Pointers to Structs and Syntactic Sugar ---
	fmt.Println("\n--- 4. Pointers to Structs & Syntactic Sugar ---")

	myCar := Car{Brand: "Toyota", Model: "Corolla", Year: 2020}
	fmt.Printf("Original Car: %+v\n", myCar)

	// Creating a pointer to a struct
	carPtr := &myCar

	// Modifying a struct through a pointer:
	// In C/C++, you'd use (*carPtr).Year = 2021 or carPtr->Year = 2021.
	// Go provides syntactic sugar: you can access fields directly through the pointer using the dot operator.
	// Go automatically dereferences the pointer under the hood.
	carPtr.Year = 2021 // Equivalent to (*carPtr).Year = 2021
	fmt.Printf("Modified via pointer: %+v\n", myCar)

	// Initializing a pointer using new()
	// new(Car) allocates memory for a Car struct, sets fields to zero values, and returns a pointer to it (*Car).
	newCarPtr := new(Car)
	newCarPtr.Brand = "Ford"
	fmt.Printf("Allocated with new(): %+v (Type: %T)\n", newCarPtr, newCarPtr)

	// --- 5. Anonymous Structs ---
	fmt.Println("\n--- 5. Anonymous Structs ---")

	// Anonymous structs are defined and initialized in a single step without a named type.
	// Great for quick, one-off structures (e.g., config data, test cases, JSON request structures).
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}
	fmt.Printf("Anonymous struct config: Host=%s, Port=%d, Type=%T\n", config.Host, config.Port, config)

	// --- 6. Embedded Structs and Field Promotion ---
	fmt.Println("\n--- 6. Embedded Structs & Field Promotion ---")

	// When a struct embeds another struct type (anonymous field), all fields of the
	// embedded struct are "promoted" to the outer struct.
	emp := Employee{
		ID:   101,
		Name: "David",
		Address: Address{
			City:    "New York",
			Country: "USA",
		},
	}

	// We can access embedded fields directly as if they were top-level fields of Employee!
	// This is called Field Promotion.
	fmt.Printf("Employee: ID=%d, Name=%s, City=%s, Country=%s\n", emp.ID, emp.Name, emp.City, emp.Country)

	// We can still access them through the embedded type name if there are name collisions:
	fmt.Printf("Access via full path: emp.Address.City = %s\n", emp.Address.City)

	// --- 7. Struct Comparison ---
	fmt.Println("\n--- 7. Struct Comparison ---")

	// Structs are comparable if all of their fields are comparable.
	// Basic types (string, int, bool) are comparable, so Person structs can be compared using == or !=.
	personA := Person{Name: "Eve", Age: 28}
	personB := Person{Name: "Eve", Age: 28}
	personC := Person{Name: "Frank", Age: 35}

	fmt.Printf("personA == personB: %t\n", personA == personB)
	fmt.Printf("personA == personC: %t\n", personA == personC)

	// NOTE: Structs containing slices, maps, or functions are NOT comparable.
	// Attempting to use == on them will result in a compile-time error.

	// --- 8. Struct Tags and JSON ---
	fmt.Println("\n--- 8. Struct Tags and JSON Serialization ---")

	user := User{
		Username: "gopher123",
		Email:    "gopher@golang.org",
		Password: "supersecretpassword",
		Website:  "", // Zero value, will be omitted because of `omitempty` tag
	}

	// Marshal struct to JSON bytes
	jsonData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Println("Marshaled JSON Output:")
	fmt.Println(string(jsonData))

	// Unmarshaling JSON back to struct
	jsonInput := `{"user_name":"dev_gopher","email":"dev@golang.org","website":"https://go.dev"}`
	var decodedUser User
	err = json.Unmarshal([]byte(jsonInput), &decodedUser)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	fmt.Printf("Decoded User Struct: %+v\n", decodedUser)
}
