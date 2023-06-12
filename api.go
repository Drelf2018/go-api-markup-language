package parser

import (
	"reflect"
	"strings"
)

// 请求任务
type Api struct {
	// 接口地址
	Url string `json:"url"`
	// 请求方式
	Method string `json:"method"`
	// 函数简介
	Hint string `json:"comment,omitempty"`
	// 函数名
	Function string `json:"-"`
	// 用来定位 info 的锚点
	Anchor string `json:"anchor,omitempty"`
	// 接口描述
	Info Tokens `json:"-"`
	// 接口载荷
	Data Tokens `json:"data,omitempty"`
	// 接口参数
	Params Tokens `json:"params,omitempty"`
	// 接口请求头
	Headers Tokens `json:"headers,omitempty"`
	// 接口文本
	Cookies Tokens `json:"cookies,omitempty"`
	// 这是一会用到的妙妙工具
	position *Tokens
}

// 为 api 添加子项
func (api *Api) Add(token *Token) {
	if token.Name == "url" {
		api.Url = token.Value
	} else if RequestTypes.Has(token.Name) {
		apiValue := reflect.ValueOf(*api)
		key := Capitalize(token.Name)
		value := apiValue.FieldByName(key).Interface().(Tokens)
		api.position = &value
	} else if token.Name == "}" {
		api.position = &api.Info
	} else {
		api.position.Add(token)
	}
}

// 构造函数
func NewApi(token *Token) *Api {
	api := Api{
		"",
		token.Type,
		token.Hint,
		token.Name,
		token.Name,
		make(Tokens),
		make(Tokens),
		make(Tokens),
		make(Tokens),
		make(Tokens),
		nil,
	}
	api.position = &api.Info
	return &api
}

// Api 管理器 用来保存和输出解析的 Api
type ApiManager struct {
	// 当前解析的任务
	*Api
	// 已经解析完成的任务
	Apis []*Api
	// 输出为 json/yml 的字典
	Output map[string]*Api
}

// 保存 api
func (am *ApiManager) Done() {
	if last := am.Api; last != nil {
		am.Apis = append(am.Apis, last)
		am.Output[last.Function] = last
		if len(last.Info) == 0 {
			last.Anchor = ""
		}
	}
}

// 添加新 api
func (am *ApiManager) New(token *Token) {
	am.Done()
	am.Api = NewApi(token)
}

// 解析为 json
func (am *ApiManager) ToJson(path string) error {
	s := JsonDump(am.Output, "    ")
	for _, api := range Filter(am.Apis, func(a *Api) bool { return a.Anchor != "" }) {
		info := JsonDump(api.Info.ToDict(), "        ")
		s = strings.Replace(s, "\"anchor\": \""+api.Function+"\"", Slice(info, "\"", "\""), 1)
	}
	return WriteFile(path, s)
}

// 构造函数
func NewApiManager() *ApiManager {
	return &ApiManager{
		nil,
		make([]*Api, 0),
		make(map[string]*Api),
	}
}
