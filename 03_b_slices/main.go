package main

import "fmt"

func main() {
	fmt.Println("==================================================")
	fmt.Println("                SLICES IN GO                      ")
	fmt.Println("==================================================")

	// 1. Declaration and Zero Values (Nil vs Empty Slices)
	// A slice is a dynamically-sized, flexible view into the elements of an array.
	// Unlike arrays, the length of a slice is not part of its type.
	fmt.Println("\n--- 1. Declaration and Zero Values ---")

	// Declaring a slice without initialization results in a nil slice
	var nilSlice []int
	fmt.Printf("var nilSlice []int -> Value: %v, len: %d, cap: %d, is nil: %t\n",
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)

	// An empty slice is initialized but contains no elements. It is NOT nil.
	emptySlice := []int{}
	fmt.Printf("emptySlice := []int{} -> Value: %v, len: %d, cap: %d, is nil: %t\n",
		emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)

	// 2. Initialization: Literals and the make() function
	fmt.Println("\n--- 2. Initialization with Literals and make() ---")

	// Slice Literal
	colors := []string{"Red", "Green", "Blue"}
	fmt.Printf("Slice literal: %v, Type: %T, len: %d, cap: %d\n",
		colors, colors, len(colors), cap(colors))

	// Creating a slice using make()
	// make([]T, length, capacity)
	// If capacity is omitted, it defaults to the length.
	makeSlice1 := make([]int, 3) // length 3, capacity 3 (elements zero-valued)
	fmt.Printf("make([]int, 3) -> Value: %v, len: %d, cap: %d\n",
		makeSlice1, len(makeSlice1), cap(makeSlice1))

	makeSlice2 := make([]int, 3, 10) // length 3, capacity 10
	fmt.Printf("make([]int, 3, 10) -> Value: %v, len: %d, cap: %d\n",
		makeSlice2, len(makeSlice2), cap(makeSlice2))

	// 3. Slicing an Existing Array or Slice
	fmt.Println("\n--- 3. Slicing Arrays and Slices ---")
	arr := [5]int{10, 20, 30, 40, 50}

	// Slice expression: array[low:high]
	// Includes index 'low', excludes index 'high'. Length is (high - low).
	s1 := arr[1:4] // elements at indices 1, 2, 3 -> [20, 30, 40]
	fmt.Printf("arr[1:4] -> Value: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	// Default indices:
	s2 := arr[:3] // same as arr[0:3] -> [10, 20, 30]
	s3 := arr[2:] // same as arr[2:len(arr)] -> [30, 40, 50]
	s4 := arr[:]  // same as arr[0:len(arr)] -> [10, 20, 30, 40, 50]
	fmt.Printf("arr[:3] -> %v\n", s2)
	fmt.Printf("arr[2:] -> %v\n", s3)
	fmt.Printf("arr[:]  -> %v\n", s4)

	// 4. Modifying Elements and Shared Memory
	// Slices do not store any data; they are just descriptors/headers pointing to an underlying array.
	// Modifying elements of a slice modifies the underlying array, and thus affects all other slices pointing to the same array!
	fmt.Println("\n--- 4. Modifying Elements (Shared Underlying Array) ---")
	sharedArray := [4]int{1, 2, 3, 4}
	sliceA := sharedArray[0:3] // [1, 2, 3]
	sliceB := sharedArray[1:4] // [2, 3, 4]

	fmt.Printf("Before modification - Array: %v, sliceA: %v, sliceB: %v\n", sharedArray, sliceA, sliceB)

	// Modify element in sliceA
	sliceA[1] = 999 // modifies index 1 of the underlying array (which is shared by sliceB at index 0)
	fmt.Printf("After modifying sliceA[1] = 999:\n")
	fmt.Printf("  Array:  %v\n", sharedArray)
	fmt.Printf("  sliceA: %v\n", sliceA)
	fmt.Printf("  sliceB: %v\n", sliceB)

	// 5. Slice Headers Internals (Pointer, Len, Cap)
	// A slice value is a small struct (slice header) containing:
	// - Pointer to the first element in the underlying array
	// - Length (number of elements in the slice)
	// - Capacity (number of elements in the underlying array, starting from the slice's first element)
	fmt.Println("\n--- 5. Slice Headers Internals (Pointers, Len, Cap) ---")
	baseArray := [6]int{10, 20, 30, 40, 50, 60}
	subSlice1 := baseArray[1:4] // elements: [20, 30, 40]
	subSlice2 := baseArray[2:5] // elements: [30, 40, 50]

	fmt.Printf("baseArray:     Pointer: %p, Len: %d, Cap: %d\n", &baseArray, len(baseArray), cap(baseArray))
	fmt.Printf("subSlice1[1:4]: Pointer: %p (starts at 20), Len: %d, Cap: %d\n", subSlice1, len(subSlice1), cap(subSlice1))
	fmt.Printf("subSlice2[2:5]: Pointer: %p (starts at 30), Len: %d, Cap: %d\n", subSlice2, len(subSlice2), cap(subSlice2))

	// 6. Growing Slices with append()
	// The append() function adds elements to the end of a slice.
	// If the capacity of the underlying array is sufficient, append reuses the underlying array.
	// If capacity is exceeded, append allocates a NEW underlying array, copies the data, and returns a slice pointing to the new array.
	fmt.Println("\n--- 6. Growing Slices with append() ---")
	originalSlice := make([]int, 0, 3) // length: 0, capacity: 3
	fmt.Printf("Initial: %v, len: %d, cap: %d, pointer: %p\n", originalSlice, len(originalSlice), cap(originalSlice), originalSlice)

	// Append elements within capacity
	sAppend1 := append(originalSlice, 1)
	sAppend2 := append(sAppend1, 2, 3) // can append multiple elements
	fmt.Printf("Appended 1, 2, 3 (within capacity):\n")
	fmt.Printf("  sAppend2: %v, len: %d, cap: %d, pointer: %p\n", sAppend2, len(sAppend2), cap(sAppend2), sAppend2)

	// Append exceeding capacity
	sAppend3 := append(sAppend2, 4) // Capacity exceeded! Go allocates a new, larger array.
	fmt.Printf("Appended 4 (exceeded capacity):\n")
	fmt.Printf("  sAppend3: %v, len: %d, cap: %d, pointer: %p (changed!)\n", sAppend3, len(sAppend3), cap(sAppend3), sAppend3)

	// Original underlying array remains unchanged (sAppend2 still points to it)
	fmt.Printf("  original sAppend2 is still: %v at pointer %p\n", sAppend2, sAppend2)

	// Appending one slice to another using the variadic operator '...'
	sliceOne := []int{100, 200}
	sliceTwo := []int{300, 400}
	combined := append(sliceOne, sliceTwo...)
	fmt.Printf("Combined sliceOne + sliceTwo: %v\n", combined)

	// 7. Copying Slices with copy()
	// The built-in copy() function copies elements from a source slice to a destination slice.
	// It returns the number of elements copied, which is the minimum of len(src) and len(dest).
	// The copy function copies to a separate underlying array, preventing shared memory side effects.
	fmt.Println("\n--- 7. Copying Slices with copy() ---")
	src := []int{1, 2, 3, 4, 5}

	// Note: copy to a slice with length 0 copies nothing!
	var destNil []int
	numCopiedNil := copy(destNil, src)
	fmt.Printf("copying to nil/empty slice: copied %d elements, destNil: %v\n", numCopiedNil, destNil)

	// Correct way: pre-allocate destination slice with correct length
	dest := make([]int, len(src))
	numCopied := copy(dest, src)
	fmt.Printf("copying to allocated slice: copied %d elements, dest: %v\n", numCopied, dest)

	// Modifying dest does not affect src (independent underlying arrays)
	dest[0] = 999
	fmt.Printf("After modifying dest[0] = 999:\n")
	fmt.Printf("  src : %v\n", src)
	fmt.Printf("  dest: %v\n", dest)

	// 8. Passing Slices to Functions
	// Slices are passed to functions by value (the slice header is copied).
	// However, because the slice header contains a pointer to the underlying array,
	// changes to slice elements inside a function ARE visible to the caller.
	fmt.Println("\n--- 8. Passing Slices to Functions ---")
	mySlice := []int{1, 2, 3}
	fmt.Printf("Before modifyElements: %v\n", mySlice)
	modifyElements(mySlice)
	fmt.Printf("After modifyElements : %v (Changed!)\n", mySlice)

	// However, actions that modify the slice header (like append) will NOT affect
	// the caller's slice header because the length and capacity copies are discarded.
	fmt.Printf("\nBefore attemptAppend: %v, len: %d, cap: %d\n", mySlice, len(mySlice), cap(mySlice))
	attemptAppend(mySlice)
	fmt.Printf("After attemptAppend : %v, len: %d, cap: %d (Unchanged!)\n", mySlice, len(mySlice), cap(mySlice))

	// To update the caller's slice structure, either return the new slice or pass a pointer.
	mySlice = correctAppend(mySlice)
	fmt.Printf("After correctAppend : %v (Updated! returned value assigned back)\n", mySlice)
}

// modifyElements changes elements of the slice. Since the slice points to the
// same underlying array, this modification is visible to the caller.
func modifyElements(s []int) {
	if len(s) > 0 {
		s[0] = 99
	}
}

// attemptAppend appends an element inside the function.
// This modifies the copy of the slice header, but the caller's slice header remains unchanged.
func attemptAppend(s []int) {
	s = append(s, 42)
	fmt.Printf("  Inside attemptAppend: %v, len: %d, cap: %d\n", s, len(s), cap(s))
}

// correctAppend returns the modified slice, which the caller must assign back.
func correctAppend(s []int) []int {
	return append(s, 42)
}
