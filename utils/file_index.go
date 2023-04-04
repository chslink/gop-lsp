package utils

import "strings"

// FileIndexInfo 文件索引信息，文件名映射到全路径
type FileIndexInfo struct {
	fileNameMap    map[string](map[string]string) // 文件名（包括后缀名）映射到所有的完整路径
	freFileNameMap map[string](map[string]string) // 文件名（不包括后缀名）映射到所有的完整路径
}

// CreateFileIndexInfo 创建文件索引对象
func CreateFileIndexInfo() *FileIndexInfo {
	return &FileIndexInfo{
		fileNameMap:    map[string](map[string]string){},
		freFileNameMap: map[string](map[string]string){},
	}
}

// InsertOneFile 缓存中插入一个文件名
func (f *FileIndexInfo) InsertOneFile(strFile string) {
	// CompleteFilePathToPreStr
	strVec := strings.Split(strFile, "/")
	fileName := strVec[len(strVec)-1]

	completePathToPreStr := CompleteFilePathToPreStr(strFile)
	if valeMap, ok := f.fileNameMap[fileName]; ok {
		valeMap[strFile] = completePathToPreStr
	} else {
		f.fileNameMap[fileName] = map[string]string{}
		f.fileNameMap[fileName][strFile] = completePathToPreStr
	}

	seperateIndex := strings.Index(fileName, ".")
	if seperateIndex < 0 {
		return
	}

	preStr := fileName[0:seperateIndex]
	if valeMap, ok := f.freFileNameMap[preStr]; ok {
		valeMap[strFile] = completePathToPreStr
	} else {
		f.freFileNameMap[preStr] = map[string]string{}
		f.freFileNameMap[preStr][strFile] = completePathToPreStr
	}
}

// RemoveOneFile 清除指定的文件
func (f *FileIndexInfo) RemoveOneFile(strFile string) {
	strVec := strings.Split(strFile, "/")
	fileName := strVec[len(strVec)-1]
	if valeMap, ok := f.fileNameMap[fileName]; ok {
		delete(valeMap, fileName)
	}

	seperateIndex := strings.Index(fileName, ".")
	if seperateIndex < 0 {
		return
	}

	preStr := fileName[0:seperateIndex]
	if valeMap, ok := f.freFileNameMap[preStr]; ok {
		delete(valeMap, preStr)
	}
}

// GetFileNameMap 获取文件名（包括后缀名）映射的所有文件名称
func (f *FileIndexInfo) GetFileNameMap(strFile string) map[string]string {
	return f.fileNameMap[strFile]
}

// GetPreFileNameMap 获取文件名（不包括后缀名）映射的所有文件名称
func (f *FileIndexInfo) GetPreFileNameMap(strFile string) map[string]string {
	return f.freFileNameMap[strFile]
}
