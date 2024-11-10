// Package protocol provides functionalities related to communication protocols.
//
// Contains the implementation of RESP2 (Redis Serialization Protocol version 2),
// which is a simple, line-based protocol for communication between Redis clients and servers.
// It defines various reply types that conform to the RESP2 specification and provides methods
// to serialize these replies into byte sequences suitable for network transmission.
package protocol

import (
	"bytes"
	"strconv"
)

var (

	// CRLF is the line separator of redis serialization protocol
	CRLF = "\r\n"
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
	return []byte("$" + strconv.Itoa(len(r.Bytes)) + CRLF + string(r.Bytes) + CRLF)
}

// ArrayReply Clients send commands to the Redis server as RESP arrays. Similarly, some Redis commands that return collections of elements use arrays as their replies.
// An example is the LRANGE command that returns elements of a list.
//
//	base format:: *<number-of-elements>\r\n<element-1>...<element-n>
type ArrayReply struct {
	BytesArr [][]byte
}

func NewArrayReply(bytesArr [][]byte) *ArrayReply {
	return &ArrayReply{
		BytesArr: bytesArr,
	}
}

func (r *ArrayReply) ToBytes() []byte {
	var buf bytes.Buffer
	bytesArrLen := len(r.BytesArr)
	// * + number-of-elements + "\r\n"
	bufSize := 1 + len(strconv.Itoa(bytesArrLen)) + 2
	for _, bs := range r.BytesArr {
		if len(bs) == 0 {
			// $-1\r\n
			bufSize += 3 + 2
		} else {
			// $<length>\r\n<data>\r\n
			bufSize += 1 + len(strconv.Itoa(len(bs))) + 2 + len(bs) + 2
		}
	}
	buf.Grow(bufSize)
	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(bytesArrLen))
	buf.WriteString(CRLF)
	for _, bs := range r.BytesArr {
		if len(bs) == 0 {
			buf.WriteString("$-1")
			buf.WriteString(CRLF)
		} else {
			buf.WriteString("$")
			buf.WriteString(strconv.Itoa(len(bs)))
			buf.WriteString(CRLF)
			buf.Write(bs)
			buf.WriteString(CRLF)
		}
	}
	return buf.Bytes()
}
