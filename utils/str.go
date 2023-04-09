package utils

import (
	"regexp"
	"strings"
)

// GetCompleteLineStr 获取一行完整的内容
func GetCompleteLineStr(contents []byte, offset int) (lineStr string) {
	conLen := len(contents)
	if offset == conLen {
		offset = offset - 1
	}
	beforeLinePos := offset
	for index := offset - 1; index >= 0; index-- {
		ch := contents[index]
		if ch == '\r' || ch == '\n' {
			break
		}
		beforeLinePos = index
	}
	endLinePos := offset
	for index := offset; index < conLen; index++ {
		ch := contents[index]
		if ch == '\r' || ch == '\n' {
			break
		}
		endLinePos = index
	}
	lineStr = string(contents[beforeLinePos : endLinePos+1])
	return lineStr
}

// GetOpenFileStr 判断是否为打开的文件，返回文件名 import
func GetOpenFileStr(contents []byte, offset int, character int, referFiles []string) []string {
	// 获取当前行的所有内容
	lineContents := GetCompleteLineStr(contents, offset)

	// 1) 引入其他文件的正则
	regDofile := regexp.MustCompile(`dofile *?\( *?\"[0-9a-zA-Z_/\-]+.lua\" *?\)`)
	regRequire := regexp.MustCompile(`require *?(\()? *?[\"|\'][0-9a-zA-Z_/\-|.]+[\"|\'] *?(\))?`)

	// ""内的内容
	regFen := regexp.MustCompile(`[\"|\'][0-9a-zA-Z_/\.\-]+[\"|\']`)

	// 是否需要.lua后缀
	needLuaSuffix := false
	requireFlag := false

	// 匹配的表达式
	importVec := regDofile.FindAllString(lineContents, -1)
	if len(importVec) == 0 {
		importVec = regRequire.FindAllString(lineContents, -1)
		if len(importVec) > 0 {
			requireFlag = true
		}
	} else {
		needLuaSuffix = true
	}

	if len(importVec) == 0 {
		for _, strOne := range referFiles {
			regImport1 := regexp.MustCompile(strOne + ` *?(\()? *?[\"|\'][0-9a-zA-Z_/|.\-]+.lua+[\"|\'] *?(\))?`)
			importVec = regImport1.FindAllString(lineContents, -1)
			if len(importVec) > 0 {
				needLuaSuffix = true
				break
			}

			regImport2 := regexp.MustCompile(strOne + ` *?(\()? *?[\"|\'][0-9a-zA-Z_/|.\-]+[\"|\'] *?(\))?`)
			importVec = regImport2.FindAllString(lineContents, -1)
			if len(importVec) > 0 {
				needLuaSuffix = false
				break
			}
		}
	}

	strOpenFile := ""
	for _, importStr := range importVec {
		findIndex := strings.Index(lineContents, importStr)
		if findIndex == -1 {
			continue
		}

		regStrFen := regFen.FindAllString(importStr, -1)
		if len(regStrFen) == 0 {
			continue
		}

		strFileTemp := regStrFen[0]
		if len(strFileTemp) < 2 {
			continue
		}

		strFileTemp = strFileTemp[1 : len(strFileTemp)-1]
		findTempIndex := strings.Index(importStr, strFileTemp)
		if findTempIndex == -1 {
			continue
		}

		importBeginIndex := findIndex + findTempIndex
		importEndIndex := findIndex + findTempIndex + len(strFileTemp)

		if character >= importBeginIndex && character <= importEndIndex {
			if needLuaSuffix && strings.HasSuffix(strFileTemp, ".lua") {
				strFileTemp = strFileTemp[0 : len(strFileTemp)-4]
			}

			strFileTemp = strings.Replace(strFileTemp, ".", "/", -1)
			strOpenFile = strFileTemp
			break
		}
	}

	if strOpenFile == "" {
		return make([]string, 0)
	}

	strModName := strings.TrimSuffix(strOpenFile, ".lua")
	result := []string{strModName + ".lua", strModName + ".so"}
	if requireFlag {
		result = append(result, strModName+"/init.lua")
	}
	return result
}
