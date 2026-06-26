package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("                  MAPS IN GO                      ")
	fmt.Println("==================================================")

	// 1. Declaration and Zero Values (Nil vs Empty Maps)
	// A map is a built-in hash table type that maps keys to values.
	// The zero value of a map is nil. A nil map has no keys, nor can keys be added.
	fmt.Println("\n--- 1. Declaration and Zero Values ---")

	var nilMap map[string]int
	fmt.Printf("var nilMap map[string]int -> Value: %v, len: %d, is nil: %t\n",
		nilMap, len(nilMap), nilMap == nil)

	// WARNING: Writing to a nil map causes a runtime panic!
	// nilMap["apple"] = 5 // Uncommenting this will panic: assignment to entry in nil map

	// An empty map is initialized and ready to use. It is NOT nil.
	emptyMap := map[string]int{}
	fmt.Printf("emptyMap := map[string]int{} -> Value: %v, len: %d, is nil: %t\n",
		emptyMap, len(emptyMap), emptyMap == nil)

	// 2. Initialization: Literals and make()
	fmt.Println("\n--- 2. Initialization with Literals and make() ---")

	// Map Literal
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
	}
	fmt.Printf("Map literal: %v, Type: %T, len: %d\n", ages, ages, len(ages))

	// Creating a map using make()
	// make(map[KeyType]ValueType)
	makeMap1 := make(map[string]int)
	fmt.Printf("make(map[string]int) -> Value: %v, len: %d\n", makeMap1, len(makeMap1))

	// Creating a map with an initial capacity hint using make()
	// This allocates memory for the specified number of elements upfront,
	// which is more efficient if you know the approximate size beforehand.
	makeMap2 := make(map[string]int, 10)
	fmt.Printf("make(map[string]int, 10) -> Value: %v, len: %d\n", makeMap2, len(makeMap2))

	// 3. Basic Operations: Insert, Update, and Retrieve
	fmt.Println("\n--- 3. Basic Operations: Insert, Update, Retrieve ---")
	m := make(map[string]string)

	// Insert elements
	m["ES"] = "Spain"
	m["FR"] = "France"
	m["IT"] = "Italy"
	fmt.Printf("After inserts: %v\n", m)

	// Update an existing key
	m["ES"] = "Espana"
	fmt.Printf("After update (ES): %v\n", m)

	// Retrieve an existing key
	country := m["FR"]
	fmt.Printf("Retrieving m[\"FR\"] -> %s\n", country)

	// Retrieve a non-existent key
	// Go returns the zero value of the value type (empty string in this case)
	missing := m["DE"]
	fmt.Printf("Retrieving missing m[\"DE\"] -> %q (zero value)\n", missing)

	// 4. Checking Key Existence: The "Comma OK" Idiom
	fmt.Println("\n--- 4. Checking Existence (Comma OK Idiom) ---")

	// To distinguish between a key that has a zero value vs one that doesn't exist at all:
	m["US"] = "" // explicit empty string

	// Checking "US" (exists, but value is empty string)
	valUS, okUS := m["US"]
	fmt.Printf("m[\"US\"] -> value: %q, exists: %t\n", valUS, okUS)

	// Checking "DE" (does not exist)
	valDE, okDE := m["DE"]
	fmt.Printf("m[\"DE\"] -> value: %q, exists: %t\n", valDE, okDE)

	// Commonly used in if-initializers:
	if name, ok := m["IT"]; ok {
		fmt.Printf("Found IT: %s\n", name)
	} else {
		fmt.Println("IT not found")
	}

	// 5. Deleting Elements
	fmt.Println("\n--- 5. Deleting Elements ---")
	shoppingList := map[string]int{"milk": 2, "bread": 1, "butter": 3}
	fmt.Printf("Before delete: %v\n", shoppingList)

	// The delete() function removes the key-value pair.
	delete(shoppingList, "bread")
	fmt.Printf("After delete(\"bread\"): %v\n", shoppingList)

	// Deleting a non-existent key is a safe, silent no-op.
	delete(shoppingList, "eggs")
	fmt.Printf("After delete(\"eggs\") (non-existent): %v\n", shoppingList)

	// 6. Iteration: Map Order is Randomized
	fmt.Println("\n--- 6. Iterating Over Maps (Randomized Order) ---")
	randomMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	// Range loop yields key and value.
	// Go intentionally randomizes map iteration order to prevent code from relying on stable ordering!
	fmt.Print("Iteration 1: ")
	for k, v := range randomMap {
		fmt.Printf("%s:%d ", k, v)
	}
	fmt.Println()

	fmt.Print("Iteration 2: ")
	for k, v := range randomMap {
		fmt.Printf("%s:%d ", k, v)
	}
	fmt.Println()

	// How to iterate in a stable, sorted order:
	// 1. Collect all keys.
	// 2. Sort the keys.
	// 3. Loop over sorted keys to get values.
	fmt.Println("Sorting and iterating:")
	keys := make([]string, 0, len(randomMap))
	for k := range randomMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, randomMap[k])
	}

	// 7. Maps are References / Passed to Functions
	// Maps are reference types. When you assign a map or pass it to a function,
	// you are copying the pointer to the underlying hmap header.
	// Thus, modifications within a function are visible to the caller.
	fmt.Println("\n--- 7. Maps passed to Functions (Reference Semantics) ---")
	userScores := map[string]int{"Alice": 95, "Bob": 80}
	fmt.Printf("Before function call: %v\n", userScores)

	modifyMap(userScores)
	fmt.Printf("After function call : %v (Changed!)\n", userScores)

	// 8. Thread Safety & Concurrency
	// Maps are NOT safe for concurrent read and write operations.
	// Simultaneous writes or a write during a read will cause a fatal runtime crash.
	fmt.Println("\n--- 8. Concurrency & Thread Safety ---")
	fmt.Println("Using sync.Mutex to protect map writes:")

	safeCounter := SafeCounter{
		m: make(map[string]int),
	}

	var wg sync.WaitGroup
	// Start multiple goroutines to increment key counter
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				safeCounter.Increment(fmt.Sprintf("worker-%d", workerID))
			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("Final safe counter counts for 5 workers (each ran 100 times):\n")
	for k, v := range safeCounter.m {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

// modifyMap modifies the map elements.
func modifyMap(m map[string]int) {
	m["Alice"] = 100 // modifies the same underlying data structure
	m["Charlie"] = 75
}

// SafeCounter wraps a map and a mutex to make it safe for concurrent use.
type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

// Increment increments the counter for a given key.
func (c *SafeCounter) Increment(key string) {
	c.mu.Lock()         // Lock the mutex before writing to the map
	defer c.mu.Unlock() // Unlock after writing
	c.m[key]++
}
