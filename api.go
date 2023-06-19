package parser

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
	Info *Token `json:"-" yaml:"-"`
	// 接口载荷
	Data *Token `json:"data,omitempty" yaml:"data,omitempty"`
	// 接口参数
	Params *Token `json:"params,omitempty" yaml:"params,omitempty"`
	// 接口请求头
	Headers *Token `json:"headers,omitempty" yaml:"headers,omitempty"`
	// 接口文本
	Cookies *Token `json:"cookies,omitempty" yaml:"cookies,omitempty"`
	// 接口返回
	Response *Token `json:"response,omitempty" yaml:"response,omitempty"`
}

// 构造函数
func NewApi(token *Token) *Api {
	return &Api{
		token.Pop("url").Value,
		token.Type,
		token.Hint,
		token.Name,
		token,
		token.Pop("data"),
		token.Pop("params"),
		token.Pop("headers"),
		token.Pop("cookies"),
		token.Pop("response"),
	}
}

// Api 管理器 用来保存和输出解析的 Api
//
// Apis: 已经解析完成的任务
//
// Output: 用来输出 json/yml 的字典
//
// Vartypes: 支持的变量类型 auto str num bool
type ApiManager struct {
	Apis     []*Api
	Output   map[string]*Api
	VarTypes *Types
}

// 添加新 api
func (am *ApiManager) Add(token *Token) {
	api := NewApi(token)
	am.Apis = append(am.Apis, api)
	am.Output[api.Function] = api
	if len(api.Info.Tokens) == 0 {
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
		NewTypes("type", "enum", "auto", "str", "num", "bool"),
	}
}
