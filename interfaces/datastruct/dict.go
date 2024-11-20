package datastruct

// Dict interface defines a simple key-value storage structure and its operation methods.
// It allows users to store and retrieve values associated with string keys.
type Dict interface {
	// Get retrieves the value associated with the given key.
	// It accepts a string key and returns the value along with a boolean indicating if the key exists.
	// If the key exists, it returns the corresponding value and true; otherwise, it returns nil and false.
	Get(key string) (val interface{}, exists bool)

	// Len returns the number of key-value pairs currently stored in the structure.
	Len() int

	// Put inserts a key-value pair into the storage structure.
	// If the key already exists, it overwrites the old value and returns the old value.
	// Parameters:
	//   key - The string key.
	//   val - The value of any type.
	// Returns:
	//   result - The type of the old value that was overwritten.
	Put(key string, val interface{}) (result int)

	// PutIfAbsent attempts to insert a key-value pair into the cache only if the key does not already exist.
	// Parameters:
	//   key - The key to insert.
	//   val - The corresponding value.
	// Returns:
	//   result - Indicates the result of the insertion operation, possibly used to indicate status or the number of entries affected.
	PutIfAbsent(key string, val interface{}) (result int)

	// PutIfExists inserts a key-value pair into the cache only if the key already exists.
	// Parameters:
	//   key - The key to insert.
	//   val - The corresponding value.
	// Returns:
	//   result - Indicates the result of the insertion operation, possibly used to indicate status or the number of entries affected.
	PutIfExists(key string, val interface{}) (result int)

	// Keys returns a list of all keys currently stored in the structure.
	// Returns:
	//   A slice of strings containing all the keys.
	Keys() []string

	// Remove removes the key-value pair associated with the specified key from the storage structure.
	// Parameters:
	//   key - The string key.
	// Returns:
	//   val - The removed value.
	//   result - Indicates the result of the removal operation.
	Remove(key string) (val interface{}, result int)

	// Clear empties the storage structure by removing all key-value pairs.
	Clear()

	// Foreach iterates over all key-value pairs in the storage structure and executes the given function f for each pair.
	// Parameters:
	//   consumer - A function that takes a key and value as parameters and returns a boolean.
	//              If f returns false, the iteration is terminated early.
	Foreach(consumer func(key string, val interface{}) bool)

	// RandomKeys generates and returns a specified number of random keys.
	// Parameters:
	//   limit - The number of random keys to generate.
	// Returns:
	//   A slice of strings containing the generated random keys.
	RandomKeys(limit int) []string

	// RandomDistinctKeys generates and returns a specified number of unique random keys.
	// Parameters:
	//   limit - The number of unique random keys to generate.
	// Returns:
	//   A slice of strings containing the generated unique random keys.
	// Note:
	//   Ensures that the returned keys are unique, suitable for scenarios requiring unique keys.
	RandomDistinctKeys(limit int) []string

	// DictScan scans the dictionary and returns keys matching the specified pattern along with the next cursor.
	// Parameters:
	//   cursor - The starting cursor for the scan.
	//   count - The number of elements to return per scan.
	//   pattern - The pattern to filter keys.
	// Returns:
	//   A slice of byte slices containing the keys that match the pattern.
	//   An integer representing the next cursor for the scan.
	// Note:
	//   Supports incremental scanning, suitable for efficient traversal of large datasets.
	DictScan(cursor int, count int, pattern string) ([][]byte, int)
}
