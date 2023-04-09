package check

import (
	"os"
	"path/filepath"
	"strings"

	"gop-lsp/logger"
	"gop-lsp/utils"
)

// DirManager 所有的目录管理
type DirManager struct {
	// 插件前端传入的工程的根目录
	vSRootDir string

	// 插件前端传入的所有次级目录，目前VScode支持多目录文件夹
	subDirVec []string

	// luahelper.json文件包含的相对路径，如果没有luahelper.json，该值默认为./
	configRelativeDir string

	// 为主工程真实的目录，为VSRootDir与ConfigRelativeDir值的拼接组合
	mainDir string
}

// CreateDirManager  default dir manager
func CreateDirManager() *DirManager {
	dirManager := &DirManager{
		vSRootDir:         "",
		subDirVec:         []string{},
		configRelativeDir: "./",
		mainDir:           "",
	}
	return dirManager
}

// InitMainDir 初始化
func (d *DirManager) InitMainDir() {
	vsRootDir := d.vSRootDir
	if vsRootDir == "" {
		return
	}

	strMainDir := ""
	if strings.HasPrefix(d.configRelativeDir, ".") {
		if d.configRelativeDir != "./" {
			strMainDir = vsRootDir + d.configRelativeDir
		} else {
			strMainDir = vsRootDir
		}
	} else {
		strMainDir = d.configRelativeDir
	}

	strMainDir, _ = filepath.Abs(strMainDir)
	strMainDir = utils.GeConvertPathFormat(strMainDir)
	strMainDir = strings.Replace(strMainDir, "./", "/", -1)

	logger.Debugf("vsRootDir=%s, strMainDir=%s", vsRootDir, strMainDir)

	d.mainDir = strMainDir
}

// SetVSRootDir set VSCode root dir
func (d *DirManager) SetVSRootDir(vsRootDir string) {
	d.vSRootDir = vsRootDir
}

// GetMainDir 获取主的目录
func (d *DirManager) GetMainDir() string {
	return d.mainDir
}

// GetVsRootDir 获取主的目录
func (d *DirManager) GetVsRootDir() string {
	return d.vSRootDir
}

// PushOneSubDir 新增一个子文件夹
func (d *DirManager) PushOneSubDir(subDir string) {
	d.subDirVec = append(d.subDirVec, subDir)
}

// RemoveOneSubDir 移除一个子文件夹
func (d *DirManager) RemoveOneSubDir(subDir string) bool {
	// 获取当前文件夹在subDirs 下的索引
	fileIndex := -1
	for index, dir := range d.subDirVec {
		if dir == subDir {
			fileIndex = index
		}
	}

	if fileIndex == -1 {
		return false
	}

	// 若移除的当前workspace 文件夹中包含的子文件夹， 则不需要做任何处理
	if !d.IsDirExistWorkspace(subDir) {
		return false
	}
	d.subDirVec = append(d.subDirVec[0:fileIndex], d.subDirVec[fileIndex+1:]...)

	return true
}

// IsDirExistWorkspace 判断当前文件夹下是否存在于 当前项目中的mainDir 和 subDirs中
func (d *DirManager) IsDirExistWorkspace(path string) bool {
	if !IsDirExist(path) {
		return false
	}

	if d.mainDir != "" && IsSubDir(path, d.mainDir) {
		return true
	}

	for _, dir := range d.subDirVec {
		if IsSubDir(path, dir) {
			return true
		}
	}
	return false
}

// GetCompletePath 传入目录和后面文件名，拼接成完整的路径
func (d *DirManager) GetCompletePath(baseDir string, fileName string) string {
	strPath := baseDir
	if !strings.HasSuffix(baseDir, "/") {
		strPath = strPath + "/"
	}

	strPath = strPath + fileName
	return strPath
}

// IsDirExist 判断所给文件夹是否存在
func IsDirExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return s.IsDir()
}

// IsFileExist 判断所给文件夹是否存在
func IsFileExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !s.IsDir()
}

// IsSubDir 判断传入的path 是否是另一个文件夹dir 下的子文件夹
func IsSubDir(path string, dir string) bool {
	if path == dir {
		return true
	}

	pathLen := len(path)
	dirLen := len(dir)

	if pathLen >= dirLen {
		return false
	}

	if strings.HasPrefix(dir, path) && dir[pathLen] == '/' {
		return true
	}

	return false
}
