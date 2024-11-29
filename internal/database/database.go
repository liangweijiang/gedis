package database

import (
	"fmt"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"github.com/liangweijiang/gedis/internal/access/protocol"
	"github.com/liangweijiang/gedis/internal/datastruct/dict"
	"strings"
)

// dataDictSize defines the size of the concurrent dictionary used for storing data entities.
const (
	dataDictSize = 1 << 16
	// ttlDictSize defines the size of the concurrent dictionary used for storing key expiration times.
	ttlDictSize = 1 << 10
)

// Database represents the database instance, containing data, TTL, and version information.
type Database struct {
	index int
	// key -> DataEntity
	data *dict.ConcurrentDict
	// key -> expireTime (time.Time)
	ttlMap *dict.ConcurrentDict
	// key -> version(uint32)
	versionMap *dict.ConcurrentDict
}

// newDatabase creates and returns a new Database instance.
func newDatabase() *Database {
	return &Database{
		index:      0,
		data:       dict.NewConcurrentDict(dataDictSize),
		ttlMap:     dict.NewConcurrentDict(ttlDictSize),
		versionMap: dict.NewConcurrentDict(dataDictSize),
	}
}

// Exec executes the given command within a database instance.
// It takes a redis.Connection and a command line (cmdLine) as parameters and returns a redis.Reply.
func (db *Database) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {
	//todo  transaction control commands and other commands which cannot execute within transaction
	return db.execNormalCommand(cmdLine)
}

// execNormalCommand executes a normal command.
// It looks up the command from the cmdTable, validates the number of arguments, and executes the command.
// If the command does not exist or the number of arguments is incorrect, it returns the corresponding error information.
func (db *Database) execNormalCommand(cmdLine [][]byte) redis.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return protocol.NewSimpleErrorReply(fmt.Sprintf("ERR unknown command `%s`, with args beginning with:", cmdName))
	}
	if cmd.validateArity(cmdLine) {
		return protocol.NewArgNumErrReply(cmdName)
	}
	prepare := cmd.prepare
	writeKeys, readKeys := prepare(cmdLine[1:])
	db.addVersion(writeKeys...)
	db.RWLocks(writeKeys, readKeys)
	defer db.RWUnLocks(writeKeys, readKeys)
	return cmd.executor(db, cmdLine[1:])
}

// addVersion adds the version for the specified keys.
// It is used to manage the version of data to support transactions or other features.
func (db *Database) addVersion(keys ...string) {
	for _, key := range keys {
		versionCode := db.GetVersion(key)
		db.versionMap.Put(key, versionCode)
	}
}

// GetVersion gets the version number for a specified key.
// If the key does not exist in the version map, it returns 0.
func (db *Database) GetVersion(key string) uint32 {
	entity, ok := db.versionMap.Get(key)
	if !ok {
		return 0
	}
	return entity.(uint32)
}

// RWLocks acquires read/write locks for the specified keys.
// It helps ensure data consistency during read/write operations.
func (db *Database) RWLocks(writeKeys []string, readKeys []string) {
	db.data.RWLocks(writeKeys, readKeys)
}

// RWUnLocks releases the read/write locks for the specified keys.
// It is used to release locks acquired by RWLocks after operations are completed.
func (db *Database) RWUnLocks(writeKeys []string, readKeys []string) {
	db.data.RWUnlocks(writeKeys, readKeys)
}
