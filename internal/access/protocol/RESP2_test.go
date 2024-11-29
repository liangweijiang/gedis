package protocol // replace with the actual package name

import (
	"bytes"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"testing"
)

// TestSimpleStringReplyToBytes is a test function for the ToBytes method of SimpleStringReply.
func TestSimpleStringReplyToBytes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []byte
	}{
		{"EmptyContent", "", []byte("+\r\n")},
		{"HelloWorld", "HelloWorld", []byte("+HelloWorld\r\n")},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SimpleStringReply{
				Content: tt.content,
			}
			if got := r.ToBytes(); !bytes.Equal(got, tt.expected) {
				t.Errorf("SimpleStringReply.ToBytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewSimpleErrorReply(t *testing.T) {
	testCases := []struct {
		err      string
		expected string
	}{
		{"", ""},
		{"An error occurred", "An error occurred"},
	}

	for _, tc := range testCases {
		reply := NewSimpleErrorReply(tc.err)
		if reply.Error != tc.expected {
			t.Errorf("Expected error message to be '%s', got '%s' instead", tc.expected, reply.Error)
		}
	}
}

func TestSimpleErrorReplyToBytes(t *testing.T) {
	testCases := []struct {
		err      string
		expected []byte
	}{
		{"error1", []byte("-error1\r\n")},
		{"error2", []byte("-error2\r\n")},
		// 更多测试用例
	}

	for _, tc := range testCases {
		reply := NewSimpleErrorReply(tc.err)
		if b := reply.ToBytes(); string(b) != string(tc.expected) {
			t.Errorf("SimpleErrorReply.ToBytes() for error '%s' = %s; want %s", tc.err, b, tc.expected)
		}
	}
}

func TestIntegerReplyToBytes(t *testing.T) {
	testCases := []struct {
		input    int64
		expected []byte
	}{
		{0, []byte(":0\r\n")},
		{123, []byte(":123\r\n")},
		{-456, []byte(":-456\r\n")},
	}

	for _, tc := range testCases {
		reply := NewIntegerReply(tc.input)
		b := reply.ToBytes()
		if !equalBytes(b, tc.expected) {
			t.Errorf("NewIntegerReply(%d).ToBytes() = %v; want %v", tc.input, b, tc.expected)
		}
	}
}

func TestBulkStringReplyToBytes(t *testing.T) {
	tests := []struct {
		name      string
		bulkReply *BulkStringReply
		want      []byte
	}{
		{
			name: "Non-empty Bytes",
			bulkReply: &BulkStringReply{
				Bytes: []byte("test"),
			},
			want: []byte("$4\r\ntest\r\n"),
		},
		{
			name: "Empty Bytes",
			bulkReply: &BulkStringReply{
				Bytes: []byte(""),
			},
			want: nullBulkBytes,
		},
		{
			name: "nil Bytes",
			bulkReply: &BulkStringReply{
				Bytes: nil,
			},
			want: nullBulkBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bulkReply.ToBytes(); !equalBytes(got, tt.want) {
				t.Errorf("BulkStringReply.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayReplyToBytes(t *testing.T) {
	// Define your test cases
	testCases := []struct {
		name     string
		input    *ArrayReply
		expected []byte
	}{
		{
			name: "Empty ArrayReply",
			input: &ArrayReply{
				Replys: nil,
			},
			expected: emptyArrayBytes,
		},
		{
			name: "ArrayReply with one SimpleStringReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewSimpleStringReply("OK"),
				},
			},
			expected: []byte("*1\r\n+OK\r\n"),
		},
		{
			name: "ArrayReply with one SimpleErrorReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewSimpleErrorReply("Error message"),
				},
			},
			expected: []byte("*1\r\n-Error message\r\n"),
		},
		{
			name: "ArrayReply with one IntegerReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewIntegerReply(123),
				},
			},
			expected: []byte("*1\r\n:123\r\n"),
		},
		{
			name: "ArrayReply with one BulkStringReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewBulkStringReply([]byte("data")),
				},
			},
			expected: []byte("*1\r\n$4\r\ndata\r\n"),
		},
		{
			name: "ArrayReply with multiple different replies",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewSimpleStringReply("OK"),
					NewSimpleErrorReply("Error message"),
					NewIntegerReply(123),
					NewBulkStringReply([]byte("data")),
				},
			},
			expected: []byte("*4\r\n+OK\r\n-Error message\r\n:123\r\n$4\r\ndata\r\n"),
		},
		{
			name: "ArrayReply with empty BulkStringReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewBulkStringReply(nil),
				},
			},
			expected: []byte("*1\r\n$-1\r\n"),
		},
		{
			name: "ArrayReply with multiple empty BulkStringReplies",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewBulkStringReply(nil),
					NewBulkStringReply(nil),
				},
			},
			expected: []byte("*2\r\n$-1\r\n$-1\r\n"),
		},
		{
			name: "ArrayReply with mixed lengths",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewSimpleStringReply("a"),
					NewSimpleStringReply("bb"),
					NewSimpleStringReply("ccc"),
				},
			},
			expected: []byte("*3\r\n+a\r\n+bb\r\n+ccc\r\n"),
		},
		{
			name: "ArrayReply with nested ArrayReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewSimpleStringReply("OK"),
					NewArrayReply([]redis.Reply{
						NewSimpleStringReply("Nested1"),
						NewSimpleStringReply("Nested2"),
					}),
					NewSimpleErrorReply("Error message"),
				},
			},
			expected: []byte("*3\r\n+OK\r\n*2\r\n+Nested1\r\n+Nested2\r\n-Error message\r\n"),
		},
		{
			name: "ArrayReply with deeply nested ArrayReply",
			input: &ArrayReply{
				Replys: []redis.Reply{
					NewArrayReply([]redis.Reply{
						NewSimpleStringReply("Level1"),
						NewArrayReply([]redis.Reply{
							NewSimpleStringReply("Level2"),
							NewArrayReply([]redis.Reply{
								NewSimpleStringReply("Level3"),
							}),
						}),
					}),
				},
			},
			expected: []byte("*1\r\n*2\r\n+Level1\r\n*2\r\n+Level2\r\n*1\r\n+Level3\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToBytes()
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}

// equalBytes checks if two slices of bytes are equal.
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
