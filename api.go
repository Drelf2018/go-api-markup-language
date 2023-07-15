package aml

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
	api := Api{
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
	if len(api.Info.Map) == 0 {
		api.Function = ""
	}
	return &api
}
