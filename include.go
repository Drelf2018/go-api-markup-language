package aml

import (
	"path/filepath"
	"strings"
)

type Include struct {
	Path string
	Args []string
	All  bool
}

func (inc *Include) Need(s string) bool {
	return inc.All || In(inc.Args, s)
}

func (inc *Include) ToApi(dir string) string {
	return filepath.Join(dir, inc.Path) + ".aml"
}

func NewInclude(path, items string) *Include {
	// 想加个 @ 的语法糖的
	// if utils.Startswith(ipath, "@") {
	// }
	args := strings.Split(items, ",")
	for i, a := range args {
		args[i] = strings.TrimSpace(a)
	}
	return &Include{
		strings.ReplaceAll(path, ".", "/"),
		args,
		In(args, "*"),
	}
}
