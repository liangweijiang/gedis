package database

import (
	"fmt"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"github.com/liangweijiang/gedis/internal/redis/protocol"
	"github.com/liangweijiang/gedis/lib/logger"
	"runtime/debug"
	"strings"
	"sync/atomic"
)

const (
	masterRole = iota
	slaveRole
)

type Server struct {
	dbSet []*atomic.Value
	// for replication
	role int32
}

func (s *Server) Exec(c redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = protocol.NewSimpleErrorReply("Err unknown")
		}
	}()
	cmdName := strings.ToLower(string(cmdLine[0]))
	if cmdName == "auth" {
		return Auth(c, cmdLine[1:])
	}
	if isAuthenticated(c) {
		return protocol.NewSimpleErrorReply("NOAUTH Authentication required.")
	}

	if cmdName == "ping" {
		return Ping(cmdLine[1:])
	}
	if cmdName == "info" {
		return Info(s, cmdLine[1:])
	}
	if cmdName == "dbsize" {
		return DBSize(s)
	}

	if cmdName == "slaveof" {
		if len(cmdLine) != 3 {
			return protocol.NewArgNumErrReply("slaveof")
		}
	}
	return SlaveOf(s, cmdLine[1:])
}

func (s *Server) Close() {
	//TODO implement me
	panic("implement me")
}
