package tcpserver

import (
	"bufio"
	"context"
	"net"
)

type EchoHandler struct {
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		b := msg
		_, _ = conn.Write(b)
	}
}

func (h *EchoHandler) Close() error {
	return nil
}
