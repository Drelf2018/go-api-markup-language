package aml

const (
	NUMBER     = iota + 4 // number
	STRING                // string
	COMMA                 // ,
	LBRACKET              // [
	RBRACKET              // ]
	LBRACE                // {
	RBRACE                // }
	LANGLE                // <
	RANGLE                // >
	LGROUP                // (
	RGROUP                // )
	AUTO                  // auto
	NUM                   // num
	STR                   // str
	TYPE                  // type
	REQUIRED              // required
	OPTIONAL              // optional
	DEPRECATE             // deprecate
	GET                   // get method
	POST                  // post method
	OPTION                // option method
	PUT                   // put method
	DELETE                // delete method
	HEAD                  // head method
	PATCH                 // patch method
	BOOL                  // bool
	FROM                  // from
	IMPORT                // import
	COLON                 // :
	IDENTIFIER            // identifier
	ASSIGNMENT            // =
)

// 最小词语单元
type Token struct {
	Kind  int
	Value string
}

// 判空
func (t *Token) IsNull() bool {
	return t.Kind <= 0
}

// 新建
func (t *Token) New(kind int, value string) {
	t.Kind = kind
	t.Value = value
}

// 设置
func (t *Token) Set(n *Token) *Token {
	t.Kind = n.Kind
	t.Value = n.Value
	return t
}

// 清除
func (t *Token) Reset() {
	t.Kind = -1
}

// 切换
func (t *Token) Shift(n *Token) {
	t.Set(n)
	n.Reset()
}

// 关键字
func GetKind(result string) int {
	switch result {
	case "from":
		return FROM
	case "import":
		return IMPORT
	case "auto":
		return AUTO
	case "num":
		return NUM
	case "str":
		return STR
	case "type":
		return TYPE
	case "bool":
		return BOOL
	case "GET":
		return GET
	case "POST":
		return POST
	case "OPTION":
		return OPTION
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "HEAD":
		return HEAD
	case "PATCH":
		return PATCH
	case "required":
		return REQUIRED
	case "optional":
		return OPTIONAL
	case "deprecate":
		return DEPRECATE
	default:
		return IDENTIFIER
	}
}
