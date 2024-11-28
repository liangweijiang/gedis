package database

import "github.com/liangweijiang/gedis/interfaces/redis"

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// DB is the interface for redis style storage engine
type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	Close()
}
