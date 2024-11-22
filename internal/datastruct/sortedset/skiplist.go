package sortedset

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"strings"
)

const (
	maxLevel = 16
)

type element struct {
	Score  float64
	Member string
}

type node struct {
	element
	backward *node
	levels   []*levelNode
}

type levelNode struct {
	forward *node
	// span 记录到下一个节点跨度
	span int64
}

type SkipList struct {
	head   *node
	tail   *node
	level  int
	length int64
}

func newSkipListNode(score float64, member string, level int) *node {
	n := &node{
		element: element{
			Score:  score,
			Member: member,
		},
		levels: make([]*levelNode, level),
	}
	for i := 0; i < level; i++ {
		n.levels[i] = new(levelNode)
	}
	return n
}

func NewSkipList() *SkipList {
	return &SkipList{
		head:   newSkipListNode(math.MinInt32, "", maxLevel),
		level:  1,
		length: 0,
	}
}

func randomLevel() int {
	total := uint64(1)<<uint64(maxLevel) - 1
	k := rand.Uint64() % total
	return maxLevel - bits.Len64(k+1) + 1
}

func (sl *SkipList) Insert(member string, score float64) {
	update := make([]*node, maxLevel)
	// rank 记录头节点到插入节点的跨度
	rank := make([]int64, maxLevel)

	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for cur.levels[i].forward != nil && (cur.levels[i].forward.Score < score ||
			(cur.levels[i].forward.Score == score && cur.levels[i].forward.Member < member)) {
			rank[i] += cur.levels[i].span
			cur = cur.levels[i].forward
		}
		update[i] = cur
	}
	fmt.Println("level: ", sl.level)
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].levels[i].span = sl.length
		}
		sl.level = level
	}
	newNode := newSkipListNode(score, member, level)
	for i := 0; i < level; i++ {
		newNode.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = newNode
		newNode.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = (rank[0] - rank[i]) + 1
	}
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] != sl.head {
		newNode.backward = update[0]
	}
	if newNode.levels[0].forward != nil {
		newNode.levels[0].forward.backward = newNode
	} else {
		sl.tail = newNode
	}
	sl.length++
}

func (sl *SkipList) Remove(member string, score float64) bool {
	update := make([]*node, maxLevel)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.levels[i].forward != nil && (cur.levels[i].forward.Score < score ||
			(cur.levels[i].forward.Score == score && cur.levels[i].forward.Member < member)) {
			cur = cur.levels[i].forward
		}
		update[i] = cur
	}
	cur = cur.levels[0].forward
	if cur != nil && cur.Member == member && cur.Score == score {
		sl.removeNode(cur, update)
		return true
	}
	return false
}

func (sl *SkipList) removeNode(cur *node, update []*node) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == cur {
			update[i].levels[i].forward = cur.levels[i].forward
			update[i].levels[i].span += cur.levels[i].span - 1
		} else {
			update[i].levels[i].span--
		}
	}
	if cur.levels[0].forward != nil {
		cur.levels[0].forward.backward = cur.backward
	} else {
		sl.tail = cur.backward
	}
	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *SkipList) String() string {
	var sb strings.Builder
	for i := sl.level - 1; i >= 0; i-- {
		cur := sl.head
		for cur != nil {
			sb.WriteString(fmt.Sprintf("(sorce:%.2f, val:%v, span:%v) ", cur.Score, cur.Member, cur.levels[i].span))
			cur = cur.levels[i].forward
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
