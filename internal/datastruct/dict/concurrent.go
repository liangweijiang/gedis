package dict

import (
	"github.com/liangweijiang/gedis/lib/utils"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentDict is a concurrent dictionary structure that supports high-concurrency read and write operations by dividing the dictionary into multiple shards.
// Reference: https://developer.aliyun.com/article/1417528
type ConcurrentDict struct {
	shards     []*shardMap
	count      int32
	shardCount uint32
}

// shardMap represents a shard of the dictionary, containing a map and a read-write lock for concurrent control.
type shardMap struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

// computeCapacity calculates the capacity of the shard, ensuring it is not less than 16 and is a power of 2.
func computeCapacity(param int) int {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

// NewConcurrentDict creates a new ConcurrentDict instance with the specified number of shards.
func NewConcurrentDict(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	shards := make([]*shardMap, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shardMap{m: make(map[string]interface{})}
	}
	return &ConcurrentDict{
		shards:     shards,
		count:      0,
		shardCount: uint32(shardCount),
	}
}

// getShard selects the shard where the key belongs based on the hash value of the key.
func (dict *ConcurrentDict) getShard(key string) *shardMap {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := utils.Fnv32(key)
	index := (dict.shardCount - 1) & hashCode
	return dict.shards[index]
}

// Get retrieves the value associated with the key and indicates whether it exists.
func (dict *ConcurrentDict) Get(key string) (val interface{}, exists bool) {
	shard := dict.getShard(key)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	val, exists = shard.m[key]
	return
}

// Len returns the total number of elements in the dictionary.
func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("dict is nil")
	}
	return int(atomic.LoadInt32(&dict.count))
}

// Put inserts a key-value pair into the dictionary, overwriting the value if the key already exists.
func (dict *ConcurrentDict) Put(key string, val interface{}) (result int) {
	shard := dict.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; !ok {
		shard.m[key] = val
		return 0
	}
	shard.m[key] = val
	dict.addCount()
	return 1
}

// PutIfAbsent inserts a key-value pair into the dictionary only if the key does not already exist.
func (dict *ConcurrentDict) PutIfAbsent(key string, val interface{}) (result int) {
	shard := dict.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; !ok {
		return 0
	}
	shard.m[key] = val
	return 1
}

// PutIfExists inserts a key-value pair into the dictionary only if the key already exists.
func (dict *ConcurrentDict) PutIfExists(key string, val interface{}) (result int) {
	shard := dict.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; !ok {
		shard.m[key] = val
		return 1
	}
	return 0
}

// Keys returns all keys in the dictionary.
func (dict *ConcurrentDict) Keys() []string {
	if dict == nil {
		panic("dict is nil")
	}
	keys := make([]string, dict.Len())
	i := 0
	dict.Foreach(func(key string, val interface{}) bool {
		if i < len(keys) {
			keys[i] = key
			i++
		} else {
			keys = append(keys, key)
		}
		return true
	})
	return keys
}

// Remove deletes the value associated with the key from the dictionary and returns the deleted value and result.
func (dict *ConcurrentDict) Remove(key string) (val interface{}, result int) {
	shard := dict.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	val, ok := shard.m[key]
	if ok {
		delete(shard.m, key)
		dict.decreaseCount()
		return val, 1
	}
	return nil, 0
}

// Clear clears all data in the dictionary.
func (dict *ConcurrentDict) Clear() {
	*dict = *NewConcurrentDict(int(dict.shardCount))
}

// Foreach iterates over each key-value pair in the dictionary, stopping the iteration if the callback function returns false.
func (dict *ConcurrentDict) Foreach(consumer func(key string, val interface{}) bool) {
	if dict == nil {
		panic("dict is nil")
	}
	for _, shard := range dict.shards {
		shard.mutex.RLock()
		f := func() bool {
			defer shard.mutex.RUnlock()
			for k, v := range shard.m {
				if !consumer(k, v) {
					return false
				}
			}
			return true
		}
		if !f() {
			break
		}
	}
}

// RandomKeys returns a specified number of random keys from the dictionary. If the limit is greater than or equal to the number of elements in the dictionary, it returns all keys.
func (dict *ConcurrentDict) RandomKeys(limit int) []string {
	if dict == nil {
		panic("dict is nil")
	}
	if limit >= int(dict.count) {
		return dict.Keys()
	}
	result := make([]string, limit)
	randR := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < limit; {
		shard := dict.shards[randR.Intn(int(dict.shardCount))]
		key := shard.RandomKey()
		if key != "" {
			result[i] = key
			i++
		}
	}
	return result
}

// RandomDistinctKeys returns a specified number of distinct random keys from the dictionary. If the limit is greater than or equal to the number of elements in the dictionary, it returns all keys.
func (dict *ConcurrentDict) RandomDistinctKeys(limit int) []string {
	if dict == nil {
		panic("dict is nil")
	}
	if limit >= int(dict.count) {
		return dict.Keys()
	}
	result := make([]string, limit)
	distinctKeys := make(map[string]struct{})
	randR := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < limit; {
		shard := dict.shards[randR.Intn(int(dict.shardCount))]
		key := shard.RandomKey()
		if key != "" {
			if _, ok := distinctKeys[key]; !ok {
				distinctKeys[key] = struct{}{}
				result[i] = key
				i++
			}

		}
	}

	return result
}

// RandomKey returns a random key from the shard. If the shard is empty, it returns an empty string.
func (shard *shardMap) RandomKey() string {
	if shard == nil {
		panic("shard is nil")
	}
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()

	for key := range shard.m {
		return key
	}
	return ""
}

// DictScan is used for incremental iteration of dictionary elements. (Not implemented)
func (dict *ConcurrentDict) DictScan(cursor int, count int, pattern string) ([][]byte, int) {
	//TODO implement me
	panic("implement me")
}

// addCount atomically increments the element count of the dictionary.
func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

// decreaseCount atomically decrements the element count of the dictionary.
func (dict *ConcurrentDict) decreaseCount() int32 {
	return atomic.AddInt32(&dict.count, -1)
}
