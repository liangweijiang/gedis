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

// Payload is a struct used to encapsulate the parsed data or error information.
type Payload struct {
	Data redis.Reply
	Err  error
}

// parse function reads data from an io.Reader and parses it into Payload objects, which are sent through the out channel.
// This function captures and handles any errors that occur during parsing.
// Parameters:
//
//	r (io.Reader): The reader from which data is read.
//	out (chan<- *Payload): The channel used to send the parsed Payload objects.
func parse(r io.Reader, out chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			// Log error information and stack trace for debugging.
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

// parseReply function reads and parses a single Redis reply from the reader.
// Parameters:
//
//	reader (*bufio.Reader): The buffered reader used to read data.
//
// Returns:
//
//	redis.Reply: The parsed Redis reply.
//	error: An error if the parsing fails.
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
		// TODO: Handle other RESP3 types
		return protocol.NewArrayReply(nil), nil
	}
}

// parseBulkString function parses a bulk string reply from the reader.
// Parameters:
//
//	reader (*bufio.Reader): The buffered reader used to read data.
//	header ([]byte): The header of the bulk string.
//
// Returns:
//
//	redis.Reply: The parsed bulk string reply.
//	error: An error if the parsing fails.
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

// parseArray function parses an array reply from the reader.
// Parameters:
//
//	reader (*bufio.Reader): The buffered reader used to read data.
//	header ([]byte): The header of the array.
//
// Returns:
//
//	redis.Reply: The parsed array reply.
//	error: An error if the parsing fails.
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

// protocolError function creates an error with a message indicating an invalid RESP protocol.
// Parameters:
//
//	msg (string): The error message.
//
// Returns:
//
//	error: The created error.
func protocolError(msg string) error {
	return errors.New("invalid RESP protocol: " + msg)
}
