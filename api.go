package parser

import (
	"encoding/json"
	"strings"

	"github.com/Drelf2020/utils"
)

type Data map[string]*Token

func (data Data) Add(token *Token) {
	data[token.Name] = token
}

func DictAdd[T Data | string](dic map[string]any, k string, v T) {
	if len(v) != 0 {
		dic[k] = v
	}
}

// 转字典
func (data Data) ToDict() map[string]any {
	dic := make(map[string]any)
	for k, v := range data {
		dic[k] = v.Value
	}
	return dic
}

// 请求任务
type Api struct {
	// 请求方式
	Method string `json:"method"`
	// 函数名
	Function string `json:"-"`
	// 函数简介
	Hint string `json:"comment"`
	// 接口
	Url string `json:"url"`
	// 对接口的描述
	Info Data `json:"info"`
	// 接口需要的参数
	Params Data `json:"params,omitempty"`
	// 接口需要的负载
	Data Data `json:"data,omitempty"`
	// 这是一会用到的妙妙工具
	position string
}

// 构造函数
func NewApi(token *Token) *Api {
	return &Api{
		token.Type,
		token.Name,
		token.Hint,
		"",
		make(Data),
		make(Data),
		make(Data),
		"info",
	}
}

// 添加子项
func (api *Api) Add(token *Token) bool {
	if token.Name == "url" {
		api.Url = token.Value
	} else if token.Name == "params" {
		api.position = "params"
	} else if token.Name == "data" {
		api.position = "data"
	} else if token.Name == "}" {
		if api.position == "info" {
			return true
		}
		api.position = "info"
	} else {
		switch api.position {
		case "params":
			api.Params.Add(token)
		case "data":
			api.Data.Add(token)
		case "info":
			api.Info.Add(token)
		}
	}
	return false
}

// Api 管理器 用来保存和输出解析的 Api
type ApiManager struct {
	// 当前解析的任务
	*Api
	// 输出为 json/yml 的字典
	Output map[string]map[string]any
	// 已经解析完成的任务
	Apis []*Api
}

// 添加新 api
func (am *ApiManager) Done(api *Api) {
	if last := am.Api; last != nil {
		dic1 := last.Info.ToDict()
		dic2 := map[string]any{
			"1url":    last.Url,
			"2method": last.Method,
		}
		DictAdd(dic2, "3comment", last.Hint)
		DictAdd(dic2, "zzparams", last.Params)
		DictAdd(dic2, "zzdata", last.Data)
		(*am).Output[last.Function] = Update(dic1, dic2)
		am.Apis = append(am.Apis, last)
	}
	am.Api = api
}

// 解析为 json
func (am *ApiManager) ToJson(path string) error {
	b, err := json.MarshalIndent(am.Output, "", "    ")
	if utils.LogErr(err) {
		return err
	}

	// 吗的这也太傻逼了 go 的 json 输出 map 只会按照 key 从小到大排序
	// 想要自定义顺序好像只能这样 别的我真不会
	s := string(b)
	s = strings.ReplaceAll(s, "1url", "url")
	s = strings.ReplaceAll(s, "2method", "method")
	s = strings.ReplaceAll(s, "3comment", "comment")
	s = strings.ReplaceAll(s, "zzparams", "params")
	s = strings.ReplaceAll(s, "zzdata", "data")

	return WriteFile(path, s)
}

// 构造函数
func NewApiManager() *ApiManager {
	return &ApiManager{
		nil,
		make(map[string]map[string]any),
		make([]*Api, 0),
	}
}
