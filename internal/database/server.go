package database

import (
	"context"
	"fmt"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"github.com/liangweijiang/gedis/internal/access/protocol"
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
	dbSet  []*atomic.Value
	ctx    context.Context
	cancel context.CancelFunc

	// configVersion stands for the version of slaveStatus config. Any change of master host/port will cause configVersion increment
	// If configVersion change has been found during slaveStatus current slaveStatus procedure will stop.
	// It is designed to abort a running slaveStatus procedure
	configVersion int32
	// for replication
	role        int32
	slaveStatus *slaveStatus
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
		return SlaveOf(s, cmdLine[1:])
	}
	// todo special commands which cannot execute within transaction

	return
}

func (s *Server) Close() {
	//TODO implement me
	panic("implement me")
}
