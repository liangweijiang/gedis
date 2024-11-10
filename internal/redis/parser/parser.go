package parser

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/liangweijiang/gedis/internal/redis/protocol"
	"github.com/liangweijiang/gedis/pkg/interfaces/redis"
	"github.com/liangweijiang/gedis/pkg/lib/logger"
	"io"
	"runtime/debug"
	"strconv"
)

type Payload struct {
	Data redis.Reply
	Err  error
}

func parse(r io.Reader, out chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err, string(debug.Stack()))
		}
	}()

	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			logger.Errorf("read line error: %+v", err)
			out <- &Payload{Err: err}
			close(out)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			continue
		}
		line = bytes.TrimSuffix(line, []byte(protocol.CRLF))
		switch line[0] {
		case '+':
			out <- &Payload{
				Data: protocol.NewSimpleStringReply(string(line[1:])),
			}
		case '-':
			out <- &Payload{
				Data: protocol.NewSimpleErrorReply(string(line[1:])),
			}
		case ':':
			value, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				protocolError(out, "illegal number "+string(line[1:]))
				continue
			}
			out <- &Payload{
				Data: protocol.NewIntegerReply(value),
			}
		}
	}
}

func protocolError(ch chan<- *Payload, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &Payload{Err: err}
}
