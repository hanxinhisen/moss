// Created by Hisen at 2022/3/3.
package gormutil

const DefaultLimit = 1000

type LimitAndOffset struct {
	Offset int
	Limit  int
}

func Unpointer(offset, limit *int64) *LimitAndOffset {
	var o, l = 0, DefaultLimit
	if offset != nil {
		o = int(*offset)
	}
	if limit != nil {
		l = int(*limit)
	}
	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
