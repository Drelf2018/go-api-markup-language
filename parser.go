package parser

import (
	"strings"

	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	// 读文本
	api := ReadFile(path)

	// 预处理 获取所有自定义类型名
	ForEach(
		VarTypes.FindTokens(api),
		func(t *Token) {
			// 解析类似 res<T1,T2> 中的子类型 T1 T2
			if text := Slice(t.Name, "<", ">", 0); text != "" {
				args := strings.Split(text, ",")
				VarTypes.Add(t.PureName(), args...)
				ForEach(args, func(s string) { VarTypes.Add(s) })
			} else {
				VarTypes.Add(t.Name)
			}
		},
		func(t *Token) bool { return t.IsType() },
	)

	// 预处理 解析所有类型 包括自定义的
	tokens := new(Tokens)
	ForEach(
		VarTypes.FindTokens(api),
		func(t *Token) {
			if t.IsType() {
				tokens = VarTypes.Get(t.PureName())
			} else if t.IsClosed() {
				tokens = nil
			} else if tokens != nil {
				tokens.Add(t)
			}
		},
		func(t *Token) bool { return t.IsClosed() || t.InType() },
	)

	// 正式解析 Api
	am = NewApiManager()
	ForEach(
		VarTypes.Union(MethodTypes).FindTokens(api),
		func(t *Token) {
			if t.IsApi() {
				am.New(t)
			} else if am.Api == nil {
				return
			} else {
				am.Add(t)
				if am.position == nil {
					am.Done()
				}
			}
		},
		func(t *Token) bool { return t.IsClosed() || t.InType() },
	)
	am.Done()
	return
}
