// Package protocol provides functionality related to communication protocols.
//
// It includes the implementation of RESP2 (Redis Serialization Protocol version 2), a simple, line-based protocol used for communication between Redis clients and servers.
// It defines various reply types conforming to the RESP2 specification and provides methods to serialize these replies into byte sequences suitable for network transmission.
package protocol

import (
	"bytes"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"strconv"
)

var (

	// CRLF is the line separator of access serialization protocol
	CRLF            = "\r\n"
	nullBulkBytes   = []byte("$-1\r\n")
	emptyArrayBytes = []byte("*0\r\n")
)

// SimpleStringReply Simple strings are encoded as a plus (+) character, followed by a string.
// The string mustn't contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
//
//	base format:: +OK\r\n
type SimpleStringReply struct {
	Content string
}

func NewSimpleStringReply(str string) *SimpleStringReply {
	return &SimpleStringReply{
		Content: str,
	}
}

func (r *SimpleStringReply) ToBytes() []byte {
	return []byte("+" + r.Content + CRLF)
}

// SimpleErrorReply RESP has specific data types for errors. Simple errors, or simply just errors, are similar to simple strings, but their first character is the minus (-) character.
// The difference between simple strings and errors in RESP is that clients should treat errors as exceptions,
// whereas the string encoded in the error type is the error message itself.
// base format: -Error message\r\n
type SimpleErrorReply struct {
	Error string
}

func NewSimpleErrorReply(err string) *SimpleErrorReply {
	return &SimpleErrorReply{
		Error: err,
	}
}

func (r *SimpleErrorReply) ToBytes() []byte {
	return []byte("-" + r.Error + CRLF)
}

// IntegerReply This type is a CRLF-terminated string that represents a signed, base-10, 64-bit integer.
//
//	base format:: :[<+|->]<value>\r\n
type IntegerReply struct {
	Number int64
}

func NewIntegerReply(num int64) *IntegerReply {
	return &IntegerReply{
		Number: num,
	}
}

func (r *IntegerReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Number, 10) + CRLF)
}

// BulkStringReply A bulk string represents a single binary string. The string can be of any size,
// but by default, Redis limits it to 512 MB (see the proto-max-bulk-len configuration directive).
//
//	base format:: $<length>\r\n<data>\r\n
type BulkStringReply struct {
	Bytes []byte
}

func NewBulkStringReply(bytes []byte) *BulkStringReply {
	return &BulkStringReply{
		Bytes: bytes,
	}
}

func (r *BulkStringReply) ToBytes() []byte {
	if r.Bytes == nil || len(r.Bytes) == 0 {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(r.Bytes)) + CRLF + string(r.Bytes) + CRLF)
}

// ArrayReply Clients send commands to the Redis server as RESP arrays. Similarly, some Redis commands that return collections of elements use arrays as their replies.
// An example is the LRANGE command that returns elements of a list.
//
//	base format:: *<number-of-elements>\r\n<element-1>...<element-n>
type ArrayReply struct {
	Replys []redis.Reply
}

func NewArrayReply(reply []redis.Reply) *ArrayReply {
	return &ArrayReply{
		Replys: reply,
	}
}

func (r *ArrayReply) ToBytes() []byte {
	if r.Replys == nil || len(r.Replys) == 0 {
		return emptyArrayBytes
	}
	var buf bytes.Buffer
	bytesArrLen := len(r.Replys)
	// * + number-of-elements + "\r\n"
	bufSize := 1 + len(strconv.Itoa(bytesArrLen)) + 2
	for _, reply := range r.Replys {
		bufSize += len(reply.ToBytes())
	}
	buf.Grow(bufSize)
	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(bytesArrLen))
	buf.WriteString(CRLF)
	for _, reply := range r.Replys {
		buf.Write(reply.ToBytes())
	}
	return buf.Bytes()
}
