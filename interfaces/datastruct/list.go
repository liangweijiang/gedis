package datastruct

// List defines an interface for a list that operates on a collection of elements of the same or different types.
type List interface {
	// Add appends a value to the end of the list.
	// Parameters:
	//   - val: The value to be added.
	// Returns:
	//   - None
	Add(val interface{})

	// Get retrieves the value at the specified index.
	// Parameters:
	//   - index: The index of the value to retrieve.
	// Returns:
	//   - The value at the specified index.
	Get(index int) (val interface{})

	// Set updates the value at the specified index.
	// Parameters:
	//   - index: The index of the value to update.
	//   - val: The new value to set.
	// Returns:
	//   - None
	Set(index int, val interface{})

	// Insert inserts a value at the specified index.
	// Parameters:
	//   - index: The index where the value should be inserted.
	//   - val: The value to insert.
	// Returns:
	//   - None
	Insert(index int, val interface{})

	// Remove removes the value at the specified index and returns it.
	// Parameters:
	//   - index: The index of the value to remove.
	// Returns:
	//   - The removed value.
	Remove(index int) (val interface{})

	// RemoveLast removes the last value in the list and returns it.
	// Parameters:
	//   - None
	// Returns:
	//   - The removed value.
	RemoveLast() (val interface{})

	// RemoveAllByVal removes all elements that match the given condition.
	// Parameters:
	//   - expected: A function that determines if an element should be removed.
	// Returns:
	//   - The number of elements removed.
	RemoveAllByVal(expected func(a interface{}) bool) int

	// RemoveByVal removes up to `count` elements that match the given condition.
	// Parameters:
	//   - expected: A function that determines if an element should be removed.
	//   - count: The maximum number of elements to remove.
	// Returns:
	//   - The number of elements removed.
	RemoveByVal(expected func(a interface{}) bool, count int) int

	// ReverseRemoveByVal removes up to `count` elements that match the given condition, starting from the end of the list.
	// Parameters:
	//   - expected: A function that determines if an element should be removed.
	//   - count: The maximum number of elements to remove.
	// Returns:
	//   - The number of elements removed.
	ReverseRemoveByVal(expected func(a interface{}) bool, count int) int

	// Len returns the number of elements in the list.
	// Parameters:
	//   - None
	// Returns:
	//   - The length of the list.
	Len() int

	// ForEach applies a function to each element in the list.
	// Parameters:
	//   - consumer: A function that takes an index and a value, and returns a boolean indicating whether to continue iteration.
	// Returns:
	//   - None
	ForEach(consumer func(i int, v interface{}) bool)

	// Contains checks if the list contains any element that matches the given condition.
	// Parameters:
	//   - expected: A function that determines if an element matches the condition.
	// Returns:
	//   - True if any element matches the condition, otherwise false.
	Contains(expected func(a interface{}) bool) bool

	// Range returns a slice of elements from the start index to the stop index (exclusive).
	// Parameters:
	//   - start: The starting index.
	//   - stop: The stopping index (exclusive).
	// Returns:
	//   - A slice of elements.
	Range(start int, stop int) []interface{}

	// ToSlice converts the list to a slice of strings.
	// Parameters:
	//   - None
	// Returns:
	//   - A slice of strings.
	ToSlice() []string
}
