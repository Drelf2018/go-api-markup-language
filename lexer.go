package parser

import (
	"strings"

	"github.com/Drelf2020/utils"
)

const (
	NUMBER     = iota + 10 // number
	STRING                 // string
	LBRACKET               // [
	RBRACKET               // ]
	LBRACE                 // {
	RBRACE                 // }
	LANGLE                 // <
	RANGLE                 // >
	LGROUP                 // (
	RGROUP                 // )
	NUM                    // num
	STR                    // str
	ENUM                   // enum
	TYPE                   // type
	QUERY                  // query
	BODY                   // body
	REQUIRED               // required
	OPTIONAL               // optional
	DEPRECATE              // deprecate
	GET                    // get method
	POST                   // post method
	OPTION                 // option method
	PUT                    // put method
	DELETE                 // delete method
	HEAD                   // head method
	PATCH                  // patch method
	BOOL                   // bool
	FROM                   // from
	IMPORT                 // import
	COLON                  // :
	IDENTIFIER             // identifier
	ASSIGNMENT             // =
)

type Token struct {
	Kind  int
	Value string
}

// 关键字
func GetKind(result string) int {
	switch result {
	case "from":
		return FROM
	case "import":
		return IMPORT
	case "num":
		return NUM
	case "str":
		return STR
	case "enum":
		return ENUM
	case "type":
		return TYPE
	case "query":
		return QUERY
	case "body":
		return BODY
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

func Parse(lexer *Scanner) chan Token {
	ch := make(chan Token)

	// 发送新 Token 到通道
	// tokens 是可选的
	// 这个函数会将 lexer 中暂存的字符（如果有）转为 Token
	// 再一并 tokens 发送至通道
	send := func(tokens ...Token) {
		if lexer.Length() != 0 {
			kind := -1
			result := lexer.Restore()
			if utils.IsNumber(result) {
				kind = NUMBER
			} else {
				kind = GetKind(result)
			}
			ch <- Token{kind, result}
		}
		utils.ForEach(tokens, func(t Token) { ch <- t })
	}

	// 延迟 异步 启动解析 Token
	defer func() {
		go func() {
			for lexer.Next() {
				s := lexer.Read()

				// 多行文本
				if lexer.HasQuotation() {
					if lexer.First() == s {
						// 多行字符串结束了
						ch <- Token{STRING, lexer.Restore()[1:]}
					} else {
						lexer.Store()
					}
					continue
				}

				// 如果开头是 # 则忽略直到换行
				if lexer.First() == "#" {
					if s != "\n" {
						continue
					}
					lexer.Restore()
				}

				// 跳过空白字符
				if strings.Contains(" \t\n\r", s) {
					send()
					continue
				}

				switch s {
				case "<":
					send(Token{LANGLE, s})
				case ">":
					send(Token{RANGLE, s})
				case "{":
					send(Token{LBRACE, s})
				case "}":
					send(Token{RBRACE, s})
				case "[":
					send(Token{LBRACKET, s})
				case "]":
					send(Token{RBRACKET, s})
				case "=":
					send(Token{ASSIGNMENT, s})
				case ":":
					send(Token{COLON, s})
				case "#":
					send()
					lexer.Store()
				default:
					lexer.Store()
				}
			}

			// 读取至文件尾 关闭文件和通道
			close(ch)
			lexer.Close()
		}()
	}()
	return ch
}
