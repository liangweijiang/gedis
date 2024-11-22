package sortedset

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSkipList(t *testing.T) {
	// Initialize a SkipList with a random seed for consistent random levels.
	rand.Seed(time.Now().UnixNano())
	skipList := NewSkipList()

	// Define test cases
	testCases := []struct {
		member string
		score  float64
	}{
		{"Alice", 1.0},
		{"Bob", 2.0},
		{"Charlie", 3.0},
		{"David", 6.0},
		{"Simple", 5.0},
		// Add more test cases as needed
	}

	// Insert members into the skip list
	for _, tc := range testCases {
		skipList.Insert(tc.member, tc.score)
		fmt.Println(skipList.String())
	}

}
