package aml

import (
	"path/filepath"
	"strings"
)

// 导入
type Include struct {
	// 被导入文件路径
	path string
	// 需要导入的参数
	args []string
	// 是否导入全部参数
	all bool
}

func (i *Include) Need(s string) bool {
	return i.all || In(i.args, s)
}

func NewInclude(dir, path string, args []string) *Include {
	// 想加个 @ 的语法糖的
	// if utils.Startswith(ipath, "@") {
	// }
	return &Include{
		filepath.Join(dir, strings.ReplaceAll(path, ".", "/")) + ".aml",
		args,
		In(args, "*"),
	}
}
