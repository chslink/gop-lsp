package utils

import (
	"unsafe"
)

// StringToBytes 实现string 转换成 []byte, 不用额外的内存分配
func StringToBytes(str string) (bytes []byte) {
	bt := *(*[]byte)(unsafe.Pointer(&str))
	return bt
}

// BytesToString 实现 []byte 转换成 string, 不需要额外的内存分配
func BytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
