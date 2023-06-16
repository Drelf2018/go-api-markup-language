package parser

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/Drelf2020/utils"
	"gopkg.in/yaml.v2"
)

// 这垃圾语言怎么连 bool 异或都没有啊
func Xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}

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
	st := strings.Index(s, start)
	sp := strings.LastIndex(s, end)
	if st == -1 || sp == -1 {
		return ""
	}
	st += (cut>>1 ^ 1) * len(start)
	sp += (cut & 1) * len(end)
	return s[st:sp]
}

// 纯净类型
func NameSlice(s string) (name string, args []string) {
	name = strings.Split(s, "<")[0]
	if text := Slice(s, "<", ">", 0); text != "" {
		depth := 0
		ForEach(
			strings.Split(text, ","),
			func(s string) {
				if depth == 0 {
					args = append(args, s)
				} else {
					args[len(args)-1] += "," + s
				}
				depth += strings.Count(s, "<") - strings.Count(s, ">")
			},
		)
	}
	return
}

// json 序列化
func JsonDump(v any, indent string) string {
	b, err := json.MarshalIndent(v, "", indent)
	utils.PanicErr(err)
	return string(b)
}

// yaml 序列化
func YamlDump(v any) string {
	b, err := yaml.Marshal(v)
	utils.PanicErr(err)
	return string(b)
}

// 过滤数组
func Filter[T any](v []T, f func(T) bool) (r []T) {
	for _, o := range v {
		if f(o) {
			r = append(r, o)
		}
	}
	return
}

// 全对
func All(expr ...bool) bool {
	for _, e := range expr {
		if !e {
			return false
		}
	}
	return true
}

// 类似 python 的 map 函数
func Map[T any, V any](f func(T) V, iter []T) (v []V) {
	for _, i := range iter {
		v = append(v, f(i))
	}
	return
}

// 条件遍历数组
//
// 好牛逼的函数
func ForEach[T any](v []T, f func(T), options ...func(T) bool) {
	Map(func(o T) int {
		if All(Map(func(t func(T) bool) bool { return t(o) }, options)...) {
			f(o)
		}
		return 0
	}, v)
}

// 复制字典
func CopyMap[T any](originalMap map[string]T) map[string]T {
	// Create the target map
	targetMap := make(map[string]T)

	// Copy from the original map to the target map
	for key, value := range originalMap {
		targetMap[key] = value
	}
	return targetMap
}

// 过滤字典
//
// 注意 这会改变字典内容 使用前请复制一份
func FilterMap[T any](v map[string]T, f func(string, T) bool) {
	for k, o := range v {
		if !f(k, o) {
			delete(v, k)
		}
	}
}

// 条件遍历字典
//
// 该方法内部会复制一份字典
func ForMap[T any](v map[string]T, f func(string, T), options ...func(string, T) bool) {
	v = CopyMap(v)
	for _, option := range options {
		FilterMap(v, option)
	}
	for k, o := range v {
		f(k, o)
	}
}

// 判断字符是否为数字
func IsNumber(v string) bool {
	// 允许一个小数点
	v = strings.Replace(v, ".", "", 1)
	_, err := strconv.Atoi(v)
	return err == nil
}

// 自动类型
func AutoType(k, v string) (typ string, val any) {
	if k != "" && k != "auto" {
		typ = k
	} else if v == "true" || v == "false" {
		typ = "bool"
	} else if v == "{" {
		typ = "dict"
	} else if IsNumber(v) {
		typ = "num"
	} else {
		typ = "str"
	}
	if v == "" {
		return
	}
	switch typ {
	case "bool":
		val = v == "true"
	case "num":
		val, _ = strconv.ParseFloat(v, 64)
	case "str":
		val = v
	}
	return
}
