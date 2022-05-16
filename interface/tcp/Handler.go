package tcp

import (
	"context"
	"net"
)

//HandleFunc 处理应用的函数
type HandleFunc func(ctx context.Context, conn net.Conn)

//Handler 基于tcp的应用服务器
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
