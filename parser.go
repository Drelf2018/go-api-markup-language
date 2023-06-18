package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

// 获取 import 语句
func GetImport(api string) [][]string {
	return regexp.MustCompile(`from +([^ ]+) +import +([\w, *]+)`).FindAllStringSubmatch(api, -1)
}

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	am = NewApiManager()
	// 读文本
	api := utils.ReadFile(path)

	// 预处理 获取 import 导入的类型
	dir := filepath.Dir(path)
	utils.ForEach(
		GetImport(api),
		func(s []string) {
			api = strings.ReplaceAll(api, s[0], "")
			ipath, types := s[1], s[2]
			// 想加个 @ 的语法糖的
			// if utils.Startswith(ipath, "@") {
			// }
			ipath = filepath.Join(dir, strings.ReplaceAll(ipath, ".", "/")) + ".aml"
			args := strings.Split(types, ",")
			for i, a := range args {
				args[i] = strings.TrimSpace(a)
			}
			utils.ForMap(
				*GetApi(ipath).VarTypes,
				func(s string, t *Token) { am.VarTypes.Add(t) },
				func(s string, t *Token) bool { return t != nil },
				func(s string, t *Token) bool { return args[0] == "*" || In(args, s) },
			)
		},
	)

	// 预处理 保存所有自定义类型名
	utils.ForEach(
		am.VarTypes.FindTokens(api),
		func(t *Token) {
			utils.ForEach(t.Args, func(s string) { am.VarTypes.Add(nil, s) })
			am.VarTypes.Add(t)
		},
		func(t *Token) bool { return t.IsType() },
	)

	// 解析 Api 以及解析所有类型 包括自定义的
	chn := ""
	token := new(Token)
	utils.ForEach(
		am.VarTypes.Union(MethodTypes).FindTokens(api),
		func(t *Token) {
			t.SetTypes(am.VarTypes)
			if t.IsType() && t.IsOpen() {
				token = am.VarTypes.Get(t.Name)
			} else if t.IsApi() {
				token = t
			} else if t.IsClose() {
				if token.IsApi() {
					am.Add(token)
				}
				token = token.Parent
			} else if chn != "" {
				token.Value += t.Name
				if t.HasQuotation(chn) {
					token.Value = utils.Slice(token.Value, chn, chn, 0)
					token = token.Parent
					chn = ""
				}
			} else if chn = t.IsMultiLine(); chn != "" {
				token.Add(t)
				token = t
			} else if token != nil {
				token.Add(t)
				if t.IsOpen() {
					token = t
				}
			}
		},
	)
	return
}
