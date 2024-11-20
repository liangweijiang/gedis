package dict

import "github.com/liangweijiang/gedis/interfaces/datastruct"

// SimpleDict is a simple dictionary implementation that satisfies the datastruct.Dict interface.
// It uses an in-memory map to store key-value pairs.
var _ datastruct.Dict = &SimpleDict{}

// SimpleDict struct definition, containing a map to store key-value pairs.
type SimpleDict struct {
	m map[string]interface{}
}

// NewSimpleDict creates and returns a new instance of SimpleDict.
func NewSimpleDict() *SimpleDict {
	return &SimpleDict{
		m: make(map[string]interface{}),
	}
}

// Put inserts a key-value pair into the dictionary. If the key already exists, it does not update the value.
// Returns 1 if the key did not exist, otherwise 0.
func (d *SimpleDict) Put(key string, val interface{}) (result int) {
	_, existed := d.m[key]
	d.m[key] = val
	if !existed {
		result = 1
	}
	return
}

// PutIfAbsent inserts a key-value pair into the dictionary only if the key does not already exist.
// Returns 1 if the insertion was successful, otherwise 0.
func (d *SimpleDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := d.m[key]
	if existed {
		return 0
	}
	d.m[key] = val
	return 1
}

// PutIfExists updates the value of a key in the dictionary only if the key already exists.
// Returns 1 if the update was successful, otherwise 0.
func (d *SimpleDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := d.m[key]
	if existed {
		d.m[key] = val
		return 1
	}
	return 0
}

// Get retrieves the value associated with a key from the dictionary.
// Returns the value and a boolean indicating whether the key exists.
func (d *SimpleDict) Get(key string) (val interface{}, exists bool) {
	val, exists = d.m[key]
	return
}

// Remove deletes a key from the dictionary and returns its value.
// Returns the value and a boolean indicating whether the key existed.
func (d *SimpleDict) Remove(key string) (val interface{}, result int) {
	val, exists := d.m[key]
	if exists {
		delete(d.m, key)
		result = 1
	}
	return
}

// Foreach iterates over all key-value pairs in the dictionary and applies the consumer function.
// If the consumer function returns false, the iteration is stopped.
func (d *SimpleDict) Foreach(consumer func(key string, val interface{}) bool) {
	for k, v := range d.m {
		if !consumer(k, v) {
			break
		}
	}
}

// Keys returns a slice containing all keys in the dictionary.
func (d *SimpleDict) Keys() []string {
	keys := make([]string, len(d.m))
	i := 0
	for k := range d.m {
		keys[i] = k
		i++
	}
	return keys
}

// RandomKeys returns a slice containing a specified number of random keys from the dictionary.
// If the specified number is greater than or equal to the total number of keys, returns all keys.
func (d *SimpleDict) RandomKeys(limit int) []string {
	if limit >= len(d.m) {
		return d.Keys()
	}
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		for k := range d.m {
			result[i] = k
			break
		}
	}
	return result
}

// RandomDistinctKeys returns a slice containing a specified number of distinct random keys from the dictionary.
// If the specified number is greater than or equal to the total number of keys, returns all keys.
func (d *SimpleDict) RandomDistinctKeys(limit int) []string {
	if limit >= len(d.m) {
		return d.Keys()
	}
	result := make([]string, limit)
	i := 0
	for k := range d.m {
		if i >= limit {
			break
		}
		result[i] = k
		i++
	}
	return result
}

// DictScan is used for incremental dictionary traversal, similar to Redis's SCAN command.
// Currently, this method is not implemented.
func (d *SimpleDict) DictScan(cursor int, count int, pattern string) ([][]byte, int) {
	//TODO implement me
	panic("implement me")
}

// Len returns the number of key-value pairs in the dictionary.
func (d *SimpleDict) Len() int {
	return len(d.m)
}

// Clear removes all key-value pairs from the dictionary.
func (d *SimpleDict) Clear() {
	d.m = make(map[string]interface{})
}
