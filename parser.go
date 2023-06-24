package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	am = NewApiManager()
	// 读文本
	api := utils.ReadFile(path)

	// 预处理 获取 import 导入的类型
	dir := filepath.Dir(path)
	utils.ForEach(
		regexp.MustCompile(`from +([^ ]+) +import +([\w, *]+)`).FindAllStringSubmatch(api, -1),
		func(s []string) {
			i := NewInclude(s)
			api = strings.ReplaceAll(api, s[0], "")
			utils.ForMap(
				*GetApi(i.ToApi(dir)).VarTypes,
				func(s string, t *Sentence) { am.VarTypes.Add(t) },
				func(s string, t *Sentence) bool { return t != nil && i.Need(s) },
			)
		},
	)

	// 预处理 保存所有自定义类型名
	utils.ForEach(
		am.VarTypes.FindSentences(api),
		func(t *Sentence) { am.VarTypes.Add(t, t.Args...) },
		func(t *Sentence) bool { return t.IsType() || t.IsEnum() },
	)

	// 解析 Api 以及解析所有类型 包括自定义的
	type Context struct {
		*Handler[*Sentence]
		parent *Sentence
		chn    string
	}
	ctx := Context{NewHandler[*Sentence](), new(Sentence), ""}
	ctx.Prepare(func(t *Sentence) {
		// 预处理
		// 父语句是列表则翻转类型和变量
		// 否则更新类型
		if ctx.parent != nil && ctx.parent.IsBracket() {
			*t = *t.Exchange(am.VarTypes)
		} else {
			t.SetTypes(am.VarTypes)
		}
	}).Add(
		// 判断是否是定义类型或定义枚举
		func(t *Sentence) bool { return (t.IsType() || t.IsEnum()) && (t.IsOpen() || t.IsBracket()) },
		func(t *Sentence) { ctx.parent = am.VarTypes.Get(t.Name) },
	).Add(
		// 判断是否是 Api
		func(t *Sentence) bool { return t.IsApi() },
		func(t *Sentence) { ctx.parent = t },
	).Add(
		// 判断是否闭合该层
		// 如果闭合的是 Api 层还要添加进 ApiManager
		func(t *Sentence) bool { return t.IsClose() },
		func(t *Sentence) {
			if ctx.parent.IsApi() {
				am.Add(ctx.parent)
			}
			ctx.parent = ctx.parent.Parent
		},
	).Add(
		// chn 不为空时把当前语句作为字符串加进上一个多行文本语句的 Value 里
		func(t *Sentence) bool { return ctx.chn != "" },
		func(t *Sentence) {
			ctx.parent.Value += t.Name
			if t.HasQuotation(ctx.chn) {
				ctx.parent.Value = utils.Slice(ctx.parent.Value, ctx.chn, ctx.chn, 0)
				ctx.parent = ctx.parent.Parent
				ctx.chn = ""
			}
		},
	).Add(
		// 判断是否是多行文本 或者有左大中括号
		func(t *Sentence) bool {
			ctx.chn = t.IsMultiLine()
			ctx.parent.Add(t)
			return ctx.chn != "" || t.IsOpen() || t.IsBracket()
		},
		func(t *Sentence) { ctx.parent = t },
	)

	utils.ForEach(am.VarTypes.Union(MethodTypes).FindSentences(api), ctx.Do)
	return
}
