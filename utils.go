package parser

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/Drelf2020/utils"
)

// 读取文件
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if utils.CheckErr(err) {
		return ""
	}
	return string(data)
}

// 写入文件
func WriteFile(path, s string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if utils.LogErr(err) {
		return err
	}

	_, err = file.WriteString(s)
	if utils.LogErr(err) {
		return err
	}
	return nil
}

var re = regexp.MustCompile(` *(?:(` + VarTypes.Join() + "|" + MethodTypes.Join() + `) )? *([^:^=^\r^\n^ ]+)(?:: *([^=^\r^\n]+))? *(?:= *([^\r^\n]+))?`)

// 找出所有语句
func FindTokens(api string, callback func(*Token)) {
	for _, s := range re.FindAllStringSubmatch(api, -1) {
		callback(NewToken(s[1:]))
	}
}

// 更新字典
func Update(dic ...map[string]any) map[string]any {
	d0 := dic[0]
	for _, d1 := range dic[1:] {
		for k, v := range d1 {
			d0[k] = v
		}
	}
	return d0
}

// 首字母大写
func Capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// 字符串切片
//
// cut 表示区间开闭
//
// 0(00) 左右都不保留
//
// 1(01) 左不保留右保留
//
// 2(10) 左保留右不保留
//
// 3(11) 左右都保留
func Slice(s, start, end string, cut int) string {
	st := strings.Index(s, start) + (cut>>1^1)*len(start)
	sp := strings.LastIndex(s, end) + (cut&1)*len(end)
	return s[st:sp]
}

// json 序列化
func JsonDump(v any, indent string) string {
	b, err := json.MarshalIndent(v, "", indent)
	utils.PanicErr(err)
	return string(b)
}

// 过滤
func Filter[T any](v []T, f func(T) bool) (r []T) {
	for _, o := range v {
		if f(o) {
			r = append(r, o)
		}
	}
	return
}
