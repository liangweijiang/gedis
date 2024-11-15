package datastruct

type Locks interface {
	Lock(key string)
	RLock(key string)
	Unlock(key string)
	RUnlock(key string)
	Locks(key ...string)
	Unlocks(key ...string)
	RLocks(key ...string)
	RUnlocks(key ...string)
	RWLocks(writeKeys []string, readKeys []string)
	RWUnlocks(writeKeys []string, readKeys []string)
}
