package b_search

import (
	"sort"
)

type Interface interface {
	Value() string
}

func Less(a, b Interface) bool {
	return a.Value() < b.Value()
}

type AInterface interface {
	sort.Interface
	Index(index int) Interface
}

func Search(i Interface, a AInterface) int {
	len := a.Len()
	mid, l, h := 0, 0, len-1
	for l <= h {
		mid = (l + h) / 2
		if Less(a.Index(mid), i) {
			l = mid + 1
		} else if Less(i, a.Index(mid)) {
			h = mid - 1
		} else {
			return mid
		}
	}
	return -1
}
