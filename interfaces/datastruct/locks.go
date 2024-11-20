package datastruct

// Locks defines a set of lock operations for coordinating access to shared resources in a distributed environment.
type Locks interface {
	// Lock acquires an exclusive lock on the specified key.
	// Parameters:
	//   key: The key to lock.
	Lock(key string)

	// RLock acquires a shared (read) lock on the specified key.
	// Parameters:
	//   key: The key to lock.
	RLock(key string)

	// Unlock releases an exclusive lock on the specified key.
	// Parameters:
	//   key: The key to unlock.
	Unlock(key string)

	// RUnlock releases a shared (read) lock on the specified key.
	// Parameters:
	//   key: The key to unlock.
	RUnlock(key string)

	// Locks acquires exclusive locks on multiple keys.
	// Parameters:
	//   key: A variadic list of keys to lock.
	Locks(key ...string)

	// Unlocks releases exclusive locks on multiple keys.
	// Parameters:
	//   key: A variadic list of keys to unlock.
	Unlocks(key ...string)

	// RLocks acquires shared (read) locks on multiple keys.
	// Parameters:
	//   key: A variadic list of keys to lock.
	RLocks(key ...string)

	// RUnlocks releases shared (read) locks on multiple keys.
	// Parameters:
	//   key: A variadic list of keys to unlock.
	RUnlocks(key ...string)

	// RWLocks acquires both exclusive and shared locks on the specified keys.
	// Parameters:
	//   writeKeys: A slice of keys to acquire exclusive locks on.
	//   readKeys: A slice of keys to acquire shared locks on.
	RWLocks(writeKeys []string, readKeys []string)

	// RWUnlocks releases both exclusive and shared locks on the specified keys.
	// Parameters:
	//   writeKeys: A slice of keys to release exclusive locks on.
	//   readKeys: A slice of keys to release shared locks on.
	RWUnlocks(writeKeys []string, readKeys []string)
}
