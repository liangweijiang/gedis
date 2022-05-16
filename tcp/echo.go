package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/liangweijiang/gedis/lib/sync/wait"

	"github.com/liangweijiang/gedis/lib/sync/atomic"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

//Close 关闭客户端链接
func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	if err := e.Conn.Close(); err != nil {
		return err
	}
	return nil
}

//EchoHandler 客户端接收到的回包，用于测试
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

//NewEchoHandler 创建EchoHandler
func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(_ context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		// 可能发生:客户端EOF，客户端超时，服务器提前关闭
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				h.activeConn.Delete(client)
			} else {

			}
			return
		}
		client.Waiting.Add(1)

		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	h.closing.Set(true)
	h.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
