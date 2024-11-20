package set

import (
	"github.com/liangweijiang/gedis/interfaces/datastruct"
	"github.com/liangweijiang/gedis/internal/datastruct/dict"
)

// Set represents a data structure that implements datastruct.Set interface.
// It is used to store unique elements.
var _ datastruct.Set = &Set{}

// Set defines a set data structure.
type Set struct {
	dict datastruct.Dict
}

// NewSet creates a new set and initializes it with the given members.
func NewSet(members ...string) *Set {
	set := &Set{
		dict: dict.NewSimpleDict(),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

// Intersect returns the intersection of multiple sets.
func Intersect(sets ...datastruct.Set) datastruct.Set {
	counts := make(map[string]int)
	newSet := NewSet()
	for _, set := range sets {
		set.ForEach(func(member string) bool {
			counts[member]++
			return true
		})
	}
	setSize := len(sets)
	for key, cnt := range counts {
		if cnt == setSize {
			newSet.Add(key)
		}
	}
	return newSet
}

// Union returns the union of multiple sets.
func Union(sets ...datastruct.Set) datastruct.Set {
	newSet := NewSet()
	for _, set := range sets {
		set.ForEach(func(member string) bool {
			newSet.Add(member)
			return true
		})
	}
	return newSet
}

// Diff returns the difference between the first set and each of the following sets.
func Diff(sets ...datastruct.Set) datastruct.Set {
	if len(sets) == 0 {
		return NewSet()
	}
	result := sets[0].ShallowCopy()
	for i := 1; i < len(sets); i++ {
		sets[i].ForEach(func(member string) bool {
			result.Remove(member)
			return true
		})
		if result.Len() == 0 {
			break
		}
	}
	return result
}

// Add inserts a new element into the set.
func (s *Set) Add(val string) int {
	if s == nil {
		panic("set is nil")
	}
	return s.dict.Put(val, nil)
}

// Remove deletes an element from the set.
func (s *Set) Remove(val string) int {
	if s == nil {
		panic("set is nil")
	}
	_, result := s.dict.Remove(val)
	return result
}

// Has checks if an element exists in the set.
func (s *Set) Has(val string) bool {
	if s == nil {
		panic("set is nil")
	}
	_, exists := s.dict.Get(val)
	return exists
}

// Len returns the number of elements in the set.
func (s *Set) Len() int {
	if s == nil {
		panic("set is nil")
	}
	return s.dict.Len()
}

// ToSlice converts the set into a slice.
func (s *Set) ToSlice() []string {
	if s == nil {
		return nil
	}
	return s.dict.Keys()
}

// ForEach iterates over each element in the set and performs the given operation.
func (s *Set) ForEach(consumer func(val string) bool) {
	if s == nil || s.dict == nil {
		return
	}
	for _, val := range s.ToSlice() {
		if !consumer(val) {
			break
		}
	}
}

// ShallowCopy creates a shallow copy of the set.
func (s *Set) ShallowCopy() datastruct.Set {
	return NewSet(s.ToSlice()...)
}

// RandomMembers returns a specified number of random members from the set.
func (s *Set) RandomMembers(limit int) []string {
	if s == nil || s.dict == nil {
		return nil
	}
	return s.dict.RandomKeys(limit)
}

// RandomDistinctMembers returns a specified number of distinct random members from the set.
func (s *Set) RandomDistinctMembers(limit int) []string {
	if s == nil || s.dict == nil {
		return nil
	}
	return s.dict.RandomDistinctKeys(limit)
}

// SetScan is used to incrementally iterate over elements in the set.
// This function is not yet implemented.
func (s *Set) SetScan(cursor int, count int, pattern string) ([][]byte, int) {
	//TODO implement me
	panic("implement me")
}
