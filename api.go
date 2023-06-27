package aml

import (
	"strings"

	"github.com/Drelf2020/utils"
)

// 请求任务
type Api struct {
	// 接口地址
	Url string `json:"url" yaml:"url"`
	// 请求方式
	Method string `json:"method" yaml:"method"`
	// 函数简介
	Hint string `json:"comment,omitempty" yaml:"comment,omitempty"`
	// 函数名
	Function string `json:"function,omitempty" yaml:"function,omitempty"`
	// 接口描述
	Info *Sentence `json:"-" yaml:"-"`
	// 接口载荷
	Data *Sentence `json:"data,omitempty" yaml:"data,omitempty"`
	// 接口参数
	Params *Sentence `json:"params,omitempty" yaml:"params,omitempty"`
	// 接口请求头
	Headers *Sentence `json:"headers,omitempty" yaml:"headers,omitempty"`
	// 接口文本
	Cookies *Sentence `json:"cookies,omitempty" yaml:"cookies,omitempty"`
	// 接口返回
	Response *Sentence `json:"response,omitempty" yaml:"response,omitempty"`
}

// 构造函数
func NewApi(sentence *Sentence) *Api {
	return &Api{
		sentence.Pop("url").Value,
		sentence.Type,
		sentence.Hint,
		sentence.Name,
		sentence,
		sentence.Pop("data"),
		sentence.Pop("params"),
		sentence.Pop("headers"),
		sentence.Pop("cookies"),
		sentence.Pop("response"),
	}
}

type Types map[string]*Sentence

// Api 管理器 用来保存和输出解析的 Api
//
// Apis: 已经解析完成的任务
//
// Output: 用来输出 json/yml 的字典
//
// Vartypes: 支持的变量类型 auto str num bool
type ApiManager struct {
	Apis   []*Api
	Output map[string]*Api
	Types
}

// 添加新 api
func (am *ApiManager) Add(sentence *Sentence) {
	utils.ForMap(
		sentence.Map,
		func(s1 string, s2 *Sentence) { s2.Use(am.Types) },
		func(s1 string, s2 *Sentence) bool { return GetKind(s1) == IDENTIFIER },
	)
	api := NewApi(sentence)
	am.Apis = append(am.Apis, api)
	am.Output[api.Function] = api
	if len(api.Info.Map) == 0 {
		api.Function = ""
	}
}

// 解析为 json
func (am *ApiManager) ToJson(path string) error {
	output := JsonDump(am.Output, "    ")
	utils.ForMap(
		am.Output,
		func(s string, a *Api) {
			info := JsonDump(a.Info.ToDict(), "        ")
			output = strings.Replace(output, "\"function\": \""+s+"\"", utils.Slice(info, "\"", "\"", 3), 1)
		},
		func(s string, a *Api) bool { return a.Function != "" },
	)
	return utils.WriteFile(path, output)
}

// 解析为 yml
func (am *ApiManager) ToYaml(path string) error {
	output := YamlDump(am.Output)
	utils.ForMap(
		am.Output,
		func(s string, a *Api) {
			info := YamlDump(map[string]map[string]string{"info": a.Info.ToDict()})
			output = strings.Replace(output, "  function: "+s+"\n", strings.Replace(info, "info:\n", "", 1), 1)
		},
		func(s string, a *Api) bool { return a.Function != "" },
	)
	return utils.WriteFile(path, output)
}

// 构造函数
func NewApiManager() *ApiManager {
	return &ApiManager{
		make([]*Api, 0),
		make(map[string]*Api),
		make(Types),
	}
}
