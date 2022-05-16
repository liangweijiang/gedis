package parser

import (
	"bufio"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/liangweijiang/gedis/lib/logger"

	"github.com/liangweijiang/gedis/interface/redis"
)

// Payload stores redis.Reply or error
type Payload struct {
	Data redis.Reply
	Err  error
}

// ParseStream 通过 io.Reader 读取数据并将结果通过 channel 将结果返回给调用者
// 流式处理的接口适合供客户端/服务端使用
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		// RESP 是以行为单位的 (RESP 是一个二进制安全的文本协议，工作于 TCP 协议上)
		// 因为行分为简单字符串和二进制安全的BulkString，我们需要封装一个 readLine 函数来兼容
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				// 遇到 IO 错误，停止读取
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}
	}
	fmt.Println(msg)
}

func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var err error
	var msg []byte
	if state.bulkLen == 0 {
		// 读取简单字符串
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, fmt.Errorf("protocol error: %s", string(msg))
		}
	} else {
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] == '\r' || msg[len(msg)-1] == '\n' {
			return nil, false, fmt.Errorf("protocol error: %s", string(msg))
		}
	}
	return msg, false, nil
}
