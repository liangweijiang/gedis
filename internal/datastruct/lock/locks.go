package lock

import (
	"github.com/liangweijiang/gedis/lib/utils"
	"sort"
	"sync"
)

// Locks is a struct used to manage a set of read-write mutexes for locking and unlocking operations based on keys.
type Locks struct {
	mutexs []*sync.RWMutex
}

// NewLocks creates and returns a new Locks instance with the specified number of read-write mutexes.
// Parameters:
//
//	size: The number of read-write mutexes to create.
//
// Returns:
//
//	A pointer to a new Locks instance.
func NewLocks(size int) *Locks {
	locks := &Locks{
		mutexs: make([]*sync.RWMutex, size),
	}
	for i := 0; i < size; i++ {
		locks.mutexs[i] = new(sync.RWMutex)
	}
	return locks
}

// getIndex calculates the corresponding mutex index based on the given key.
// Parameters:
//
//	key: The key to calculate the index for.
//
// Returns:
//
//	The index of the mutex.
func (l *Locks) getIndex(key string) uint32 {
	hashCode := utils.Fnv32(key)
	return (uint32(len(l.mutexs)) - 1) & hashCode
}

// getMutex retrieves the corresponding read-write mutex based on the given key.
// Parameters:
//
//	key: The key to retrieve the mutex for.
//
// Returns:
//
//	A pointer to the read-write mutex.
func (l *Locks) getMutex(key string) *sync.RWMutex {
	return l.mutexs[l.getIndex(key)]
}

// toLockIndices calculates the list of mutex indices based on the given keys and sorts them if required.
// Parameters:
//
//	keys: A slice of keys to calculate the indices for.
//	reverse: A boolean indicating whether to sort the indices in reverse order.
//
// Returns:
//
//	A slice of mutex indices.
func (l *Locks) toLockIndices(keys []string, reverse bool) []uint32 {
	indices := make([]uint32, 0, len(keys))
	for _, key := range keys {
		index := l.getIndex(key)
		indices = append(indices, index)
	}
	sort.Slice(indices, func(i, j int) bool {
		if reverse {
			return indices[i] > indices[j]
		} else {
			return indices[i] < indices[j]
		}
	})
	return indices
}

// Lock locks the write lock for the given key.
// Parameters:
//
//	key: The key to lock.
func (l *Locks) Lock(key string) {
	mu := l.getMutex(key)
	mu.Lock()
}

// RLock locks the read lock for the given key.
// Parameters:
//
//	key: The key to lock.
func (l *Locks) RLock(key string) {
	mu := l.getMutex(key)
	mu.RLock()
}

// Unlock unlocks the write lock for the given key.
// Parameters:
//
//	key: The key to unlock.
func (l *Locks) Unlock(key string) {
	mu := l.getMutex(key)
	mu.Unlock()
}

// RUnlock unlocks the read lock for the given key.
// Parameters:
//
//	key: The key to unlock.
func (l *Locks) RUnlock(key string) {
	mu := l.getMutex(key)
	mu.RUnlock()
}

// Locks locks multiple keys with write locks, locking in index order to avoid deadlocks.
// Parameters:
//
//	key: A variadic list of keys to lock.
func (l *Locks) Locks(key ...string) {
	indices := l.toLockIndices(key, false)
	for _, index := range indices {
		l.mutexs[index].Lock()
	}
}

// RLocks locks multiple keys with read locks, locking in index order to avoid deadlocks.
// Parameters:
//
//	key: A variadic list of keys to lock.
func (l *Locks) RLocks(key ...string) {
	indices := l.toLockIndices(key, false)
	for _, index := range indices {
		l.mutexs[index].RLock()
	}
}

// Unlocks unlocks multiple keys with write locks, unlocking in reverse index order to avoid deadlocks.
// Parameters:
//
//	key: A variadic list of keys to unlock.
func (l *Locks) Unlocks(key ...string) {
	indices := l.toLockIndices(key, true)
	for _, index := range indices {
		l.mutexs[index].Unlock()
	}
}

// RUnlocks unlocks multiple keys with read locks, unlocking in reverse index order to avoid deadlocks.
// Parameters:
//
//	key: A variadic list of keys to unlock.
func (l *Locks) RUnlocks(key ...string) {
	indices := l.toLockIndices(key, true)
	for _, index := range indices {
		l.mutexs[index].RUnlock()
	}
}

// RWLocks locks a set of write keys and read keys, applying write locks to write keys and read locks to read keys, in index order to avoid deadlocks.
// Parameters:
//
//	writeKeys: A slice of keys to lock with write locks.
//	readKeys: A slice of keys to lock with read locks.
func (l *Locks) RWLocks(writeKeys []string, readKeys []string) {
	writeKeys = append(writeKeys, readKeys...)
	indices := l.toLockIndices(writeKeys, false)
	writeIndexSet := make(map[uint32]struct{})
	for _, wKey := range writeKeys {
		index := l.getIndex(wKey)
		writeIndexSet[index] = struct{}{}
	}
	for _, index := range indices {
		if _, ok := writeIndexSet[index]; ok {
			l.mutexs[index].Lock()
		} else {
			l.mutexs[index].RLock()
		}
	}
}

// RWUnlocks unlocks a set of write keys and read keys, releasing write locks from write keys and read locks from read keys, in reverse index order to avoid deadlocks.
// Parameters:
//
//	writeKeys: A slice of keys to unlock with write locks.
//	readKeys: A slice of keys to unlock with read locks.
func (l *Locks) RWUnlocks(writeKeys []string, readKeys []string) {
	writeKeys = append(writeKeys, readKeys...)
	indices := l.toLockIndices(writeKeys, true)
	writeIndexSet := make(map[uint32]struct{})
	for _, wKey := range writeKeys {
		index := l.getIndex(wKey)
		writeIndexSet[index] = struct{}{}
	}
	for _, index := range indices {
		if _, ok := writeIndexSet[index]; ok {
			l.mutexs[index].Unlock()
		} else {
			l.mutexs[index].RUnlock()
		}
	}
}
