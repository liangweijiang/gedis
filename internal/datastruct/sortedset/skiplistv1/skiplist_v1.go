package skiplistv1

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	// MaxLevel MAX_LEVEL 最高层数
	MaxLevel = 16
)

const (
	Success = iota
	NilValue
	ValueExists
)

type skipListNode struct {
	score    float64
	val      interface{}
	level    int
	forwards []*skipListNode
}

type SkipList struct {
	head   *skipListNode
	level  int
	length int
}

func newSkipListNode(score float64, val interface{}, level int) *skipListNode {
	return &skipListNode{
		score:    score,
		val:      val,
		level:    level,
		forwards: make([]*skipListNode, level),
	}
}

func NewSkipList() *SkipList {
	return &SkipList{
		head:   newSkipListNode(math.MinInt32, nil, MaxLevel),
		level:  1,
		length: 0,
	}
}

func (sl *SkipList) Insert(val interface{}, score float64) int {
	if val == nil {
		return NilValue
	}
	update := make([]*skipListNode, MaxLevel)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].score == score {
				return ValueExists
			}
			if cur.forwards[i].score > score {
				break
			}
			cur = cur.forwards[i]
		}
		update[i] = cur
	}
	level := sl.randomLevel()
	newNode := newSkipListNode(score, val, level)
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}
	for i := 0; i < level; i++ {
		next := update[i].forwards[i]
		update[i].forwards[i] = newNode
		newNode.forwards[i] = next
	}
	sl.length++
	return Success
}

func (sl *SkipList) Find(score float64) (interface{}, bool) {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].score == score {
				return cur.forwards[i].val, true
			}
			if cur.forwards[i].score > score {
				break
			}
			cur = cur.forwards[i]
		}
	}
	return nil, false
}

func (sl *SkipList) Delete(score float64) bool {
	var deleted bool
	update := make([]*skipListNode, MaxLevel)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		update[i] = cur
		for cur.forwards[i] != nil {
			if cur.forwards[i].score == score {
				deleted = true
				// sl.length--
				break
			}
			cur = cur.forwards[i]
		}
	}
	if !deleted {
		return false
	}
	cur = update[0].forwards[0]
	for i := cur.level - 1; i >= 0; i-- {
		if update[i] == sl.head && cur.forwards[i] == nil {
			sl.level = i
		}
		if update[i].forwards[i] != nil {
			update[i].forwards[i] = update[i].forwards[i].forwards[i]
		}
	}
	sl.length--
	return deleted
}

func (sl *SkipList) Len() int {
	return sl.length
}

func (sl *SkipList) Level() int {
	return sl.level
}

func (sl *SkipList) String() string {
	var sb strings.Builder
	for i := sl.level - 1; i >= 0; i-- {
		cur := sl.head
		for cur != nil {
			sb.WriteString(fmt.Sprintf("(sorce:%.2f, val:%v) ", cur.score, cur.val))
			cur = cur.forwards[i]
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (sl *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < 0.5 && level < MaxLevel {
		level++
	}
	return level
}
