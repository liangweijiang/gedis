package list

import (
	"fmt"
)

// LinkedList represents a linked list data structure.
type LinkedList struct {
	size int
	head *listNode
	tail *listNode
}

// listNode represents a node in the linked list.
type listNode struct {
	prev *listNode
	next *listNode
	val  interface{}
}

// Add appends a new value to the end of the linked list.
// Parameters:
//
//	val: The value to append
func (l *LinkedList) Add(val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	node := &listNode{
		prev: l.tail,
		next: nil,
		val:  val,
	}
	if l.tail != nil {
		l.tail.next = node
	} else {
		l.head = node
	}
	l.tail = node
	l.size++
}

// Get retrieves the value at the specified index in the linked list.
// Parameters:
//
//	index: The index of the node
//
// Returns:
//
//	interface{}: The value at the specified index
func (l *LinkedList) Get(index int) (val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	node := l.find(index)
	return node.val
}

// Set sets the value at the specified index in the linked list.
// Parameters:
//
//	index: The index of the node
//	val: The new value to set
func (l *LinkedList) Set(index int, val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	node := l.find(index)
	node.val = val
}

// Insert inserts a new value at the specified index in the linked list.
// Parameters:
//
//	index: The index at which to insert
//	val: The value to insert
func (l *LinkedList) Insert(index int, val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	if index < 0 || index > l.size {
		panic("index out of bound")
	}
	if index == l.size {
		l.Add(val)
		return
	}
	node := l.find(index)
	newNode := &listNode{
		prev: node.prev,
		next: node,
		val:  val,
	}
	if node.prev != nil {
		node.prev.next = newNode
	} else {
		l.head = newNode
	}
	node.prev = newNode
	l.size++
}

// removeNode removes the specified node from the linked list.
// Parameters:
//
//	node: The node to remove
func (l *LinkedList) removeNode(node *listNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	l.size--
}

// Remove removes the node at the specified index from the linked list.
// Parameters:
//
//	index: The index of the node to remove
//
// Returns:
//
//	interface{}: The value of the removed node
func (l *LinkedList) Remove(index int) (val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	if index < 0 || index >= l.size {
		panic("index out of bound")
	}
	node := l.find(index)
	l.removeNode(node)
	return node.val
}

// RemoveLast removes the last node from the linked list.
// Returns:
//
//	interface{}: The value of the removed node
func (l *LinkedList) RemoveLast() (val interface{}) {
	if l == nil {
		panic("list is nil")
	}
	tail := l.tail
	l.removeNode(tail)
	return tail.val
}

// RemoveAllByVal removes all nodes that satisfy the condition function from the linked list.
// Parameters:
//
//	expected: Condition function, accepts a value and returns a boolean
//
// Returns:
//
//	int: The number of nodes removed
func (l *LinkedList) RemoveAllByVal(expected func(a interface{}) bool) int {
	if l == nil {
		panic("list is nil")
	}
	var removedCnt int
	node := l.head
	for node != nil {
		next := node.next
		if expected(node.val) {
			l.removeNode(node)
			removedCnt++
		}
		node = next
	}
	return removedCnt
}

// RemoveByVal removes the first `count` nodes that satisfy the condition function from the linked list.
// Parameters:
//
//	expected: Condition function, accepts a value and returns a boolean
//	count: The number of nodes to remove
//
// Returns:
//
//	int: The actual number of nodes removed
func (l *LinkedList) RemoveByVal(expected func(a interface{}) bool, count int) int {
	if l == nil {
		panic("list is nil")
	}
	var removedCnt int
	node := l.head
	for node != nil {
		next := node.next
		if expected(node.val) {
			l.removeNode(node)
			removedCnt++
		}
		if removedCnt == count {
			break
		}
		node = next
	}
	return removedCnt
}

// ReverseRemoveByVal removes the first `count` nodes that satisfy the condition function from the end of the linked list.
// Parameters:
//
//	expected: Condition function, accepts a value and returns a boolean
//	count: The number of nodes to remove
//
// Returns:
//
//	int: The actual number of nodes removed
func (l *LinkedList) ReverseRemoveByVal(expected func(a interface{}) bool, count int) int {
	if l == nil {
		panic("list is nil")
	}
	var removedCnt int
	node := l.tail
	for node != nil {
		prev := node.prev
		if expected(node.val) {
			l.removeNode(node)
			removedCnt++
		}
		if removedCnt == count {
			break
		}
		node = prev
	}
	return removedCnt
}

// Len returns the length of the linked list.
// Returns:
//
//	int: The length of the linked list
func (l *LinkedList) Len() int {
	if l == nil {
		panic("list is nil")
	}
	return l.size
}

// ForEach executes the given function for each node in the linked list.
// Parameters:
//
//	f: Function that accepts an index and a value and returns a boolean; if false is returned, the traversal stops
func (l *LinkedList) ForEach(consumer func(i int, v interface{}) bool) {
	if l == nil {
		panic("list is nil")
	}
	node := l.head
	var i int
	for node != nil {
		if !consumer(i, node.val) {
			break
		}
		i++
		node = node.next
	}
}

// Contains checks if there is a node that satisfies the condition in the linked list.
// Parameters:
//
//	expected: Condition function, accepts a value and returns a boolean
//
// Returns:
//
//	bool: Whether there is a node that satisfies the condition
func (l *LinkedList) Contains(expected func(a interface{}) bool) bool {
	if l == nil {
		panic("list is nil")
	}
	contains := false
	l.ForEach(func(i int, val interface{}) bool {
		if expected(val) {
			contains = true
			return false
		}
		return true
	})
	return contains
}

// Range returns the sublist from `start` to `stop` in the linked list.
// Parameters:
//
//	start: The starting index of the sublist
//	stop: The ending index of the sublist (exclusive)
//
// Returns:
//
//	[]interface{}: The sublist
func (l *LinkedList) Range(start int, stop int) []interface{} {
	if l == nil {
		panic("list is nil")
	}
	if start < 0 || start >= l.size {
		panic("`start` out of range")
	}
	if stop < start || stop > l.size {
		panic("`stop` out of range")
	}
	var result []interface{}
	node := l.head
	var i int
	for node != nil {
		if i >= start && i < stop {
			result = append(result, node.val)
		} else if i >= stop {
			break
		}
		i++
		node = node.next
	}
	return result
}

// ToSlice converts the linked list to a slice of strings.
// Returns:
//
//	[]string: The slice of strings
func (l *LinkedList) ToSlice() []string {
	if l == nil {
		return nil
	}
	node := l.head
	result := make([]string, l.size)
	i := 0
	for node != nil {
		result[i] = fmt.Sprintf("%v", node.val)
		node = node.next
		i++
	}
	return result
}

// find locates the node at the specified index.
// Parameters:
//
//	index: The index of the node to find
//
// Returns:
//
//	*listNode: The node at the specified index
func (l *LinkedList) find(index int) *listNode {
	if index < 0 || index >= l.size {
		panic("index out of bound")
	}
	node := l.head
	for i := 0; i < index; i++ {
		node = node.next
	}
	return node
}
