package parser

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	am = NewApiManager()
	api := ReadFile(path)
	FindTokens(api, func(token *Token) {
		if token.IsApi() {
			am.Done(NewApi(token))
			return
		}
		am.Add(token)
	})
	am.Done(nil)
	return
}
