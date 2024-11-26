package sortedset

import (
	"errors"
	"github.com/liangweijiang/gedis/interfaces/datastruct"
	"strconv"
)

const (
	scoreNegativeInf int8 = -1
	scorePositiveInf int8 = 1
)

// ScoreBorder is a struct represents `min` `max` parameter of redis command `ZRANGEBYSCORE`
// can accept:
//
//	int or float value, such as 2.718, 2, -2.718, -2 ...
//	exclusive int or float value, such as (2.718, (2, (-2.718, (-2 ...
//	infinity: +inf, -inf， inf(same as +inf)
type ScoreBorder struct {
	Inf     int8
	Value   float64
	Exclude bool
}

var scoreNegativeInfBorder = &ScoreBorder{Inf: scoreNegativeInf}
var scorePositiveInfBorder = &ScoreBorder{Inf: scorePositiveInf}

func ParseScoreBorder(s string) (datastruct.Border, error) {
	if s == "inf" || s == "+inf" {
		return scorePositiveInfBorder, nil
	}
	if s == "-inf" {
		return scoreNegativeInfBorder, nil
	}
	if s[0] == '(' {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, errors.New("ERR min or max is not a float")
		}
		return &ScoreBorder{
			Inf:     0,
			Value:   value,
			Exclude: true,
		}, nil
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("ERR min or max is not a float")
	}
	return &ScoreBorder{
		Inf:     0,
		Value:   value,
		Exclude: false,
	}, nil
}

func (b *ScoreBorder) Greater(e *datastruct.Element) bool {
	value := e.Score
	if b.Inf == scoreNegativeInf {
		return false
	} else if b.Inf == scorePositiveInf {
		return true
	}
	if b.Exclude {
		return b.Value > value
	}
	return b.Value >= value
}

func (b *ScoreBorder) Less(e *datastruct.Element) bool {
	value := e.Score
	if b.Inf == scoreNegativeInf {
		return true
	} else if b.Inf == scorePositiveInf {
		return false
	}
	if b.Exclude {
		return b.Value < value
	}
	return b.Value <= value
}

func (b *ScoreBorder) GetValue() interface{} {
	return b.Value
}

func (b *ScoreBorder) GetExclude() bool {
	return b.Exclude
}

func (b *ScoreBorder) IsIntersected(max datastruct.Border) bool {
	minValue := b.Value
	maxValue := max.(*ScoreBorder).Value
	return minValue > maxValue || (minValue == maxValue && (b.GetExclude() || max.GetExclude()))
}
