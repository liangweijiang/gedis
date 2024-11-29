package database

import (
	"context"
	"github.com/liangweijiang/gedis/internal/access/parser"
	"net"
	"sync"
	"time"
)

type slaveStatus struct {
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	// configVersion stands for the version of slaveStatus config. Any change of master host/port will cause configVersion increment
	// If configVersion change has been found during slaveStatus current slaveStatus procedure will stop.
	// It is designed to abort a running slaveStatus procedure
	configVersion int32

	masterHost string
	masterPort int

	masterConn   net.Conn
	masterChan   <-chan *parser.Payload
	replId       string
	replOffset   int64
	lastRecvTime time.Time
	running      sync.WaitGroup
}

func (s *slaveStatus) setupMaster() {
	// todo connect to master
}
