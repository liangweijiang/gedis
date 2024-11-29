package connection

import (
	"github.com/liangweijiang/gedis/interfaces/redis"
	"net"
	"sync"
)

var _ redis.Connection = &Connection{}

type Connection struct {
	conn     net.Conn
	password string

	// selected db
	selectedDB int
}

var connPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

func (c *Connection) Write(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	return c.conn.Write(bytes)
}

func (c *Connection) Close() error {
	_ = c.conn.Close()
	connPool.Put(c)
	return nil
}

func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) SetPassword(password string) {
	c.password = password
}

func (c *Connection) GetPassword() string {
	return c.password
}

func (c *Connection) InMultiState() bool {
	return false
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(idx int) {
	c.selectedDB = idx
}
