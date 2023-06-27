package aml

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Drelf2020/utils"
	"gopkg.in/yaml.v2"
)

// 纯净类型
func NameSlice(s string) (name string, args []string) {
	name = strings.Split(s, "<")[0]
	if text := utils.Slice(s, "<", ">", 0); text != "" {
		depth := 0
		utils.ForEach(
			strings.Split(text, ","),
			func(s string) {
				s = strings.TrimSpace(s)
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

// 自动类型
func AutoType(k, v string) (typ string, val any) {
	if k != "" && k != "auto" {
		typ = k
	} else if v == "true" || v == "false" {
		typ = "bool"
	} else if v == "{" {
		typ = "dict"
	} else if v == "[" {
		typ = "list"
	} else if utils.IsNumber(v) {
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

// python list in
func In(ls []string, v string) bool {
	return len(utils.Filter(ls, func(t string) bool { return t == v })) != 0
}
