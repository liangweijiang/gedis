package parser // replace with your package name

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	testCases := []struct {
		name       string
		input      string
		wantData   string
		wantErr    error
		payloadErr bool
	}{
		{
			name:       "Normal data",
			input:      "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			wantData:   "mockReply",
			wantErr:    nil,
			payloadErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bufio.NewReader(bytes.NewReader([]byte(tc.input)))
			outChan := make(chan *Payload, 1) // Buffered channel to avoid blocking

			go parse(r, outChan)
			payload := <-outChan
			if tc.payloadErr {
				if payload.Err == nil || payload.Err.Error() != tc.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tc.wantErr, payload.Err)
				}
			} else {
				fmt.Println(string(payload.Data.ToBytes()))
			}
		})
	}
}
