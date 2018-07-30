package util

import (
	"bytes"
	"encoding/binary"
)

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	if tmp < 0 {
		tmp = -tmp
	}
	return int(tmp)
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func MoreThanHalf(n int) int {
	half := n/2 + 1
	return half
}
