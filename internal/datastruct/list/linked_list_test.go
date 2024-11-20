package list

import (
	"fmt"
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	tests := []struct {
		name     string
		vals     []interface{}
		expected []interface{}
	}{
		{
			name: "empty list",
			vals: nil,
		},
		{
			name: "single element",
			vals: []interface{}{1},
		},
		{
			name: "multiple elements",
			vals: []interface{}{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList(tt.vals...)
			fmt.Println(list.ToSlice())
		})
	}
}

func TestLinkedList(t *testing.T) {
	vals := []interface{}{1, 2, 3}
	list := NewLinkedList(vals...)
	fmt.Println(list.ToSlice())

	list.Add(4)
	fmt.Println(list.ToSlice())

	list.Remove(2)
	fmt.Println(list.ToSlice())

	list.Insert(2, 5)
	fmt.Println(list.ToSlice())
	fmt.Println("len: ", list.Len())

	fmt.Println("get: ", list.Get(3))
}
