package x_random

import (
	"math/rand"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

func Krand(size int, kind int) string {
	if size < 1 {
		return ""
	}
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func RandomString(size int) string {
	return Krand(size, KC_RAND_KIND_ALL)
}

func RandomUpper(size int) string {
	return Krand(size, KC_RAND_KIND_UPPER)
}

func RandomLower(size int) string {
	return Krand(size, KC_RAND_KIND_LOWER)
}

func RandomNumber(size int) string {
	return Krand(size, KC_RAND_KIND_NUM)
}
