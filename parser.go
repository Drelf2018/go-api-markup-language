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
				func(s string, t *Token) { am.VarTypes.Add(t) },
				func(s string, t *Token) bool { return t != nil && i.Need(s) },
			)
		},
	)

	// 预处理 保存所有自定义类型名
	utils.ForEach(
		am.VarTypes.FindTokens(api),
		func(t *Token) { am.VarTypes.Add(t, t.Args...) },
		func(t *Token) bool { return t.IsType() || t.IsEnum() },
	)

	// 解析 Api 以及解析所有类型 包括自定义的
	type Cache struct {
		parent *Token
		chn    string
	}
	hd := NewHandler[*Token, *Cache](&Cache{new(Token), ""})

	// 预处理 更新类型
	hd.Add(func(t *Token, c *Cache) bool { t.SetTypes(am.VarTypes); return false }, nil)

	// 判断是否是定义类型或定义枚举
	hd.Add(
		func(t *Token, c *Cache) bool { return (t.IsType() || t.IsEnum()) && t.IsOpen() },
		func(t *Token, c *Cache) { c.parent = am.VarTypes.Get(t.Name) },
	)

	// 判断是否是 Api
	hd.Add(
		func(t *Token, c *Cache) bool { return t.IsApi() },
		func(t *Token, c *Cache) { c.parent = t },
	)

	// 判断是否闭合该层
	// 如果闭合的是 Api 层还要添加进 ApiManager
	hd.Add(
		func(t *Token, c *Cache) bool { return t.IsClose() },
		func(t *Token, c *Cache) { am.Add(c.parent); c.parent = c.parent.Parent },
	)

	// chn 不为空时把当前语句作为字符串加进上一个多行文本语句的 Value 里
	hd.Add(
		func(t *Token, c *Cache) bool { return c.chn != "" },
		func(t *Token, c *Cache) {
			c.parent.Value += t.Name
			if t.HasQuotation(c.chn) {
				c.parent.Value = utils.Slice(c.parent.Value, c.chn, c.chn, 0)
				c.parent = c.parent.Parent
				c.chn = ""
			}
		},
	)

	// 判断是否是多行文本 或者有左大括号
	hd.Add(
		func(t *Token, c *Cache) bool {
			c.chn = t.IsMultiLine()
			c.parent.Add(t)
			return c.chn != "" || t.IsOpen()
		},
		func(t *Token, c *Cache) { c.parent = t },
	)

	utils.ForEach(am.VarTypes.Union(MethodTypes).FindTokens(api), hd.Do)
	return
}
