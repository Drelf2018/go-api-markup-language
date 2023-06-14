package parser

// 从文件解析出 Api
func GetApi(path string) (am *ApiManager) {
	// 读文本
	api := ReadFile(path)

	// 预处理 获取所有自定义类型名
	ForEach(
		VarTypes.FindTokens(api),
		func(t *Token) {
			ForEach(t.Args, func(s string) { VarTypes.Add(nil, s) })
			VarTypes.Add(t)
		},
		func(t *Token) bool { return t.IsType() },
	)

	// 解析 Api 解析所有类型 包括自定义的
	am = NewApiManager()
	token := new(Token)
	ForEach(
		VarTypes.Union(MethodTypes).FindStrings(api),
		func(sList []string) {
			t := NewToken(sList[1:]...)
			if t.IsType() && t.IsOpen() {
				token = VarTypes.Get(t.Name)
			} else if t.IsApi() {
				token = t
			} else if t.IsClose() {
				if token.IsApi() {
					am.Add(token)
				}
				token = token.Parent
			} else if token != nil {
				token.Add(t)
				if t.IsOpen() {
					token = t
				}
			}
		},
	)
	return am
}
