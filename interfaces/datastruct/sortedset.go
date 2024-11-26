package datastruct

type Element struct {
	Score  float64
	Member string
}

type Border interface {
	Greater(e *Element) bool
	Less(e *Element) bool
	GetValue() interface{}
	GetExclude() bool
	IsIntersected(max Border) bool
}

type SortedSet interface {
	Add(member string, score float64) bool
	Len() int64
	Get(member string) (*Element, bool)
	Remove(member string) bool
	GetRank(member string, desc bool) int64
	ForEachByRank(start int64, stop int64, desc bool, consumer func(element *Element) bool)
	RangeByRank(start int64, stop int64, desc bool) []*Element
	RangeCount(min Border, max Border) int64
	ForEach(min Border, max Border, offset int64, limit int64, desc bool, consumer func(element *Element) bool)
	Range(min Border, max Border, offset int64, limit int64, desc bool) []*Element
	RemoveRange(min Border, max Border) int64
	PopMin(count int) []*Element
	RemoveByRank(start int64, stop int64) int64
	ZSetScan(cursor int, count int, pattern string) ([][]byte, int)
}
