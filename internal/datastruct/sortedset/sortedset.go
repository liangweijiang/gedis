package sortedset

import (
	"github.com/liangweijiang/gedis/interfaces/datastruct"
	"strconv"
)

type SortedSet struct {
	skipList *SkipList
	dict     map[string]*datastruct.Element
}

func (s *SortedSet) Add(member string, score float64) bool {
	element, ok := s.dict[member]
	if ok {
		if element.Score == score {
			return false
		}
		s.skipList.Remove(member, element.Score)
	}
	s.skipList.Insert(member, score)
	s.dict[member] = &datastruct.Element{
		Member: member,
		Score:  score,
	}
	return true
}

func (s *SortedSet) Len() int64 {
	return s.skipList.length
}

func (s *SortedSet) Get(member string) (element *datastruct.Element, ok bool) {
	element, ok = s.dict[member]
	return element, ok
}

func (s *SortedSet) Remove(member string) bool {
	element, ok := s.dict[member]
	if ok {
		s.skipList.Remove(member, element.Score)
		delete(s.dict, member)
		return true
	}
	return false
}

func (s *SortedSet) GetRank(member string, desc bool) (rank int64) {
	element, ok := s.dict[member]
	if !ok {
		return -1
	}
	rank = s.skipList.GetRank(member, element.Score)
	if desc {
		rank = s.skipList.length - rank
	} else {
		rank--
	}
	return rank
}

func (s *SortedSet) ForEachByRank(start int64, stop int64, desc bool, consumer func(element *datastruct.Element) bool) {
	size := s.Len()
	if start < 0 || start >= size {
		panic("illegal start " + strconv.FormatInt(start, 10))
	}
	if stop < start || stop > size {
		panic("illegal end " + strconv.FormatInt(stop, 10))
	}
	var cur *Node
	if desc {
		cur = s.skipList.tail
		if start > 0 {
			cur = s.skipList.GetByRank(size - start)
		}
	} else {
		cur = s.skipList.head
		if start > 0 {
			cur = s.skipList.GetByRank(start + 1)
		}
	}
	sliceSize := int(stop - start)
	for i := 0; i < sliceSize; i++ {
		if !consumer(&cur.Element) {
			break
		}
		if desc {
			cur = cur.backward
		} else {
			cur = cur.levels[0].forward
		}
	}
}

func (s *SortedSet) RangeByRank(start int64, stop int64, desc bool) []*datastruct.Element {
	sliceSize := int(stop - start)
	slice := make([]*datastruct.Element, sliceSize)
	i := 0
	s.ForEachByRank(start, stop, desc, func(element *datastruct.Element) bool {
		slice[i] = element
		i++
		return true
	})
	return slice
}

func (s *SortedSet) RangeCount(min datastruct.Border, max datastruct.Border) int64 {
	var count int64
	s.ForEach(min, max, 0, s.Len(), false, func(element *datastruct.Element) bool {
		if !min.Less(element) {
			return true
		}
		if !max.Greater(element) {
			return false
		}
		count++
		return true
	})
	return count
}

func (s *SortedSet) ForEach(min datastruct.Border, max datastruct.Border, offset int64, limit int64, desc bool, consumer func(element *datastruct.Element) bool) {
	var cur *Node
	if desc {
		cur = s.skipList.GetLastInRange(min, max)
	} else {
		cur = s.skipList.GetFirstInRange(min, max)
	}
	for cur != nil && offset > 0 {
		if desc {
			cur = cur.backward
		} else {
			cur = cur.levels[0].forward
		}
		offset--
	}
	for i := int64(0); (i < limit || limit < 0) && cur != nil; i++ {
		if !consumer(&cur.Element) {
			break
		}
		if desc {
			cur = cur.backward
		} else {
			cur = cur.levels[0].forward
		}
		if cur == nil {
			break
		}
		if !min.Less(&cur.Element) || !max.Greater(&cur.Element) {
			break
		}
	}
}

func (s *SortedSet) Range(min datastruct.Border, max datastruct.Border, offset int64, limit int64, desc bool) []*datastruct.Element {
	result := make([]*datastruct.Element, 0)
	if limit == 0 || offset < 0 {
		return result
	}
	s.ForEach(min, max, offset, limit, desc, func(element *datastruct.Element) bool {
		result = append(result, element)
		return true
	})
	return result
}

func (s *SortedSet) RemoveRange(min datastruct.Border, max datastruct.Border) int64 {
	removed := s.skipList.RemoveRange(min, max, 0)
	for _, element := range removed {
		delete(s.dict, element.Member)
	}
	return int64(len(removed))
}

func (s *SortedSet) PopMin(count int) []*datastruct.Element {
	first := s.skipList.GetFirstInRange(scoreNegativeInfBorder, scorePositiveInfBorder)
	border := &ScoreBorder{
		Value:   first.Score,
		Exclude: false,
	}
	removed := s.skipList.RemoveRange(border, scorePositiveInfBorder, count)
	for _, element := range removed {
		delete(s.dict, element.Member)
	}
	return removed
}

func (s *SortedSet) RemoveByRank(start int64, stop int64) int64 {
	removed := s.skipList.RemoveRangeByRank(start, stop)
	for _, element := range removed {
		delete(s.dict, element.Member)
	}
	return int64(len(removed))
}

func (s *SortedSet) ZSetScan(cursor int, count int, pattern string) ([][]byte, int) {
	//TODO implement me
	panic("implement me")
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		skipList: NewSkipList(),
		dict:     make(map[string]*datastruct.Element),
	}
}
