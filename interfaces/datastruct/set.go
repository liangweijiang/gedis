package datastruct

// Set is an interface representing a collection of elements based on a hash table.
type Set interface {

	// Add adds a member to the set.
	// It returns the number of members in the set after the addition.
	// Parameters:
	//   - val: The value to add to the set.
	// Returns:
	//   - int: The new size of the set.
	Add(val string) int

	// Remove removes a member from the set.
	// It returns the number of members in the set after the removal.
	// Parameters:
	//   - val: The value to remove from the set.
	// Returns:
	//   - int: The new size of the set.
	Remove(val string) int

	// Has checks if a value exists in the set.
	// Parameters:
	//   - val: The value to check for existence.
	// Returns:
	//   - bool: True if the value exists in the set, false otherwise.
	Has(val string) bool

	// Len returns the number of members in the set.
	// Returns:
	//   - int: The size of the set.
	Len() int

	// ToSlice converts the set to a slice of strings.
	// Returns:
	//   - []string: A slice containing all members of the set.
	ToSlice() []string

	// ForEach iterates over each member of the set.
	// Parameters:
	//   - consumer: A function that takes a member as input and returns a boolean.
	//               If the function returns false, the iteration stops.
	ForEach(consumer func(member string) bool)

	// ShallowCopy creates a shallow copy of the set.
	// Returns:
	//   - Set: A new set containing the same members as the original set.
	ShallowCopy() Set

	// RandomMembers returns a specified number of random members from the set.
	// The returned members may contain duplicates.
	// Parameters:
	//   - limit: The number of random members to return.
	// Returns:
	//   - []string: A slice of random members.
	RandomMembers(limit int) []string

	// RandomDistinctMembers returns a specified number of distinct random members from the set.
	// The returned members will not contain duplicates.
	// Parameters:
	//   - limit: The number of distinct random members to return.
	// Returns:
	//   - []string: A slice of distinct random members.
	RandomDistinctMembers(limit int) []string

	// SetScan scans the set using a cursor and a pattern.
	// Parameters:
	//   - cursor: The starting cursor for the scan.
	//   - count: The number of elements to return.
	//   - pattern: A pattern to match the members.
	// Returns:
	//   - [][]byte: A slice of byte slices containing the matched members.
	//   - int: The next cursor for the scan.
	SetScan(cursor int, count int, pattern string) ([][]byte, int)
}
