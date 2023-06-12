package parser

import "github.com/Drelf2020/utils"

var log = utils.GetLog()

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	am = NewApiManager()
	api := ReadFile(path)
	FindTokens(api, func(token *Token) {
		if token.IsApi() {
			am.New(token)
		} else {
			am.Add(token)
		}
	})
	am.Done()
	return
}
