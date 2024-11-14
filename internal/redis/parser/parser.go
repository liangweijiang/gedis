package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"github.com/liangweijiang/gedis/internal/redis/protocol"
	"github.com/liangweijiang/gedis/lib/logger"
	"io"
	"runtime/debug"
	"strconv"
)

type Payload struct {
	Data redis.Reply
	Err  error
}

// parse函数负责从io.Reader中读取数据，并解析成Payload对象，通过通道out发送出去。
// 该函数在解析过程中，会捕获并处理任何发生的错误。
// 参数:
//
//	r (io.Reader): 数据来源的读取器。
//	out (chan<- *Payload): 用于发送解析后的Payload对象的通道。
func parse(r io.Reader, out chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			// 记录错误信息和堆栈追踪，便于调试。
			logger.Error(err, string(debug.Stack()))
		}
	}()

	reader := bufio.NewReader(r)
	for {
		reply, err := parseReply(reader)
		if err != nil {
			out <- &Payload{Err: err}
			close(out)
			return
		}
		out <- &Payload{Data: reply}
	}
}

func parseReply(reader *bufio.Reader) (redis.Reply, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Errorf("read line error: %+v", err)
		return nil, err
	}
	length := len(line)
	if length <= 2 || line[length-2] != '\r' {
		return nil, protocolError("invalid RESP protocol: line does not end with \\r\\n")
	}
	line = bytes.TrimSuffix(line, []byte(protocol.CRLF))
	switch line[0] {
	case '+':
		return protocol.NewSimpleStringReply(string(line[1:])), nil
	case '-':
		return protocol.NewSimpleErrorReply(string(line[1:])), nil
	case ':':
		value, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, protocolError("illegal number " + string(line[1:]))
		}
		return protocol.NewIntegerReply(value), nil
	case '$':
		reply, err := parseBulkString(reader, line)
		if err != nil {
			return nil, err
		}
		return reply, nil

	case '*':
		reply, err := parseArray(reader, line)
		if err != nil {
			return nil, err
		}
		return reply, nil
	default:

		// todo
		//Nulls	             RESP3	Simple	    _
		//Booleans	         RESP3	Simple	    #
		//Doubles	         RESP3	Simple	    ,
		//Big numbers	     RESP3	Simple	    (
		//Bulk errors	     RESP3	Aggregate	!
		//Verbatim strings	 RESP3	Aggregate	=
		//Maps	             RESP3	Aggregate	%
		//Attributes	     RESP3	Aggregate	`
		//Sets	             RESP3	Aggregate	~
		//Pushes	         RESP3	Aggregate	>
		return protocol.NewArrayReply(nil), nil
	}
}

func parseBulkString(reader *bufio.Reader, header []byte) (redis.Reply, error) {
	strLen, err := strconv.Atoi(string(header[1:]))
	if err != nil || strLen < -1 {
		return nil, protocolError("illegal bulk string header: " + string(header))
	} else if strLen == -1 {
		return protocol.NewBulkStringReply(nil), nil
	}
	bulkLine, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	length := len(bulkLine)
	if length <= 2 || bulkLine[length-2] != '\r' {
		return nil, protocolError(fmt.Sprintf("line does not end with \\r\\n"))

	}

	if len(bulkLine) != strLen+len(protocol.CRLF) {
		return nil, protocolError(fmt.Sprintf("invalid bulk string length: %d", strLen))
	}

	return protocol.NewBulkStringReply(bulkLine[:length-2]), nil
}

func parseArray(reader *bufio.Reader, header []byte) (redis.Reply, error) {
	arrLen, err := strconv.Atoi(string(header[1:]))
	if err != nil || arrLen < 0 {
		return nil, protocolError("illegal array header: " + string(header))
	} else if arrLen == 0 {
		return protocol.NewArrayReply(nil), nil
	}
	var replys []redis.Reply
	for i := 0; i < arrLen; i++ {
		reply, err := parseReply(reader)
		if err != nil {
			return nil, err
		}
		replys = append(replys, reply)
	}
	return protocol.NewArrayReply(replys), nil
}

func protocolError(msg string) error {
	return errors.New("invalid RESP protocol: " + msg)
}
