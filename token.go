package aml

import "sync"

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

// 设置
func (t *Token) Set(kind int, value string) *Token {
	t.Kind = kind
	t.Value = value
	return t
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

// Token 池
type Pool[T any] struct {
	sync.Pool
}

func (p *Pool[T]) Get() *T {
	return p.Pool.Get().(*T)
}

func (p *Pool[T]) Put(t *T) {
	// log.Debug("put | ", t)
	p.Pool.Put(t)
}

func NewPool[T any]() *Pool[T] {
	pool := Pool[T]{}
	pool.New = func() interface{} {
		return new(T)
	}
	return &pool
}
