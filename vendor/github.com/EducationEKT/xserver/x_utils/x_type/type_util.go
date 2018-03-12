package x_type

import (
	"fmt"
	"reflect"
	"strconv"
)

func V2String(v interface{}) string {
	if reflect.TypeOf(v).String() == "[]uint8" {
		return string(v.([]uint8))
	} else if reflect.TypeOf(v).String() == "[]byte" {
		return string(v.([]byte))
	}
	return fmt.Sprintf("%v", v)
}

func GetInt64(v interface{}) (int64, bool) {
	str := V2String(v)
	value, err := strconv.Atoi(str)
	if err != nil {
		return -1, false
	}
	return int64(value), true
}

func GetInt32(v interface{}) (int32, bool) {
	str := V2String(v)
	value, err := strconv.Atoi(str)
	if err != nil {
		return -1, false
	}
	return int32(value), true
}

func GetBool(v interface{}) (bool, bool) {
	str := V2String(v)
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, false
	}
	return value, true
}

func GetFloat32(v interface{}) (float64, bool) {
	str := V2String(v)
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return -1, false
	}
	return value, true
}

func GetFloat64(v interface{}) (float64, bool) {
	str := V2String(v)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1, false
	}
	return value, true
}
