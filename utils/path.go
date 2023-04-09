package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// preFixStr 工程项目的前缀，windows平台的前缀是 file:///， mac平台的前缀是file://
var preFixStr string = "file:///"

// InitialRootURIAndPath 初始化前缀
func InitialRootURIAndPath(rootURI, rootPath string) {
	rootURI, _ = url.QueryUnescape(rootURI)
	rootURI = strings.Replace(rootURI, "\\", "/", -1)

	rootPath, _ = url.QueryUnescape(rootPath)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)

	if len(rootURI) < 8 {
		return
	}

	subRootURI := rootURI[7:]
	if subRootURI == rootPath {
		preFixStr = "file://"
	}
}

// GetRemovePreStr 如果字符串以 ./为前缀，去除掉前缀
func GetRemovePreStr(str string) string {
	if strings.HasPrefix(str, "./") {
		str = str[2:]
	}
	return str
}

// VscodeURIToString 插件传入的路径转换
// vscode 传来的路径： file:///g%3A/luaproject
// 统一转换为：g%3A/luaproject，去掉前缀的file:///，并且都是这样的/../
func VscodeURIToString(strURL string) string {
	fileURL := strings.Replace(strURL, preFixStr, "", 1)
	fileURL, _ = url.QueryUnescape(fileURL)
	fileURL = strings.Replace(string(fileURL), "\\", "/", -1)

	return fileURL
}

// StringToVscodeURI 文件真实路径转换成类似的 file:///g%3A/luaproject/src/tutorial.lua"
func StringToVscodeURI(strPath string) string {
	strPath = strings.Replace(string(strPath), "\\", "/", -1)
	strEncode := strPath
	strURI := preFixStr + strEncode
	return strURI
}

// GeConvertPathFormat 文件路径统一为
func GeConvertPathFormat(strPath string) string {
	strDir := strings.Replace(strPath, "\\", "/", -1)
	return strDir
}

// CompleteFilePathToPreStr 给定的文件名，截取前面.部分的字符串， 例如 test.lua ，返回test
func CompleteFilePathToPreStr(pathFile string) (preStr string) {
	// 完整路径提前前缀
	// 字符串中，查找第一个.
	seperateIndex := strings.Index(pathFile, ".")
	if seperateIndex < 0 {
		return ""
	}

	preStr = pathFile[0:seperateIndex]
	return preStr
}

// OffsetForPosition Previously used bytes converted to rune.
// Now use the high bit to determine how many bits the character occupies.
// posLine (zero-based, from 0)
func OffsetForPosition(contents []byte, posLine, posCh int) (int, error) {
	line := 0
	col := 0
	offset := 0

	getCharBytes := func(b byte) int {
		num := 0
		for b&(1<<uint32(7-num)) != 0 {
			num++
		}
		return num
	}

	for index := 0; index < len(contents); index++ {
		if line == posLine && col == posCh {
			return offset, nil
		}

		if (line == posLine && col > posCh) || line > posLine {
			return 0, fmt.Errorf("character %d (zero-based) is beyond line %d boundary (zero-based)", posCh, posLine)
		}

		curChar := contents[index]
		if curChar > 127 {
			curCharBytes := getCharBytes(curChar)
			index += curCharBytes - 1
			offset += curCharBytes - 1
		}
		offset++
		if curChar == '\n' {
			line++
			col = 0
		} else {
			col++
		}

	}
	if line == posLine && col == posCh {
		return offset, nil
	}
	if line == 0 {
		return 0, fmt.Errorf("character %d (zero-based) is beyond first line boundary", posCh)
	}
	return 0, fmt.Errorf("file only has %d lines", line+1)
}
