package check

import (
	"gop-lsp/utils"
)

// 整体检查的入口函数

// AllProject 所有工程包含的内容
type AllProject struct {
	// 所有需要分析的文件map。key值为前缀，截取前面.部分的字符串, 例如 test.lua, 返回test。
	allFilesMap map[string]string

	// 文件名的缓存信息
	fileIndexInfo *utils.FileIndexInfo

	// // entryFile string
	// // 所有的工程入口分析文件
	// entryFilesList []string

	// // 插件客户端路径额外扩展的Lua文件夹文件map
	// clientExpFileMap map[string]struct{}

	// // 第一阶段所有文件的结果
	// fileStructMap   map[string]*results.FileStruct
	// fileStructMutex sync.Mutex //  第一阶段所有文件的结果的互斥锁

	// // 保存vscode客户端实时输入产生无法语法错误的文件结果 *FileStruct，代码提示会优先查找这个结构，这个是最新的
	// fileLRUMap *common.LRUCache

	// // 第二阶段分析的工程
	// analysisSecondMap map[string]*results.SingleProjectResult

	// // 第三阶段分析的全局结果
	// thirdStruct *results.AnalysisThird

	// // 管理所有的注释创建的type类型，key值为名称，value是这个类型的列表，允许多个存在
	// createTypeMap map[string]common.CreateTypeList

	// // 代码补全cache
	// completeCache *common.CompleteCache

	// // 整体分析的阶段数
	// checkTerm results.CheckTerm
}

// IsNeedHandle 给一个文件名，判断是否要进行处理
func (a *AllProject) IsNeedHandle(strFile string) bool {
	// todo  判断该文件是否是忽略处理的 暂时处理所有文件

	return true
}
