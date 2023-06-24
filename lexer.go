package parser

import "strings"

const (
	NUMBER     = iota // number
	STRING            //string
	LBRACKET          // [
	RBRACKET          // ]
	LBRACE            //  {
	RBRACE            //  }
	LANGLE            // <
	RANGLE            // >
	LGROUP            // (
	RGROUP            // )
	STR               // str
	ENUM              //enum
	TYPE              //type
	QUERY             //query
	BODY              //body
	REQUIRED
	OPTIONAL
	DEPRECATE
	GET               //get method
	POST              //post method
	OPTION            //option method
	PUT               //put method
	DELETE            //delete method
	HEAD              //head method
	PATCH             //patch method
	BOOL              // bool
	FROM              // from
	IMPORT            // import
	COLON             //:
	IDENTIFIER        // identifier
	ASSIGNMENT        // =
	EOF               // end
)

type LToken struct {
	Kind  int
	Value string
}

func (c *LToken) New(kind int, value string) LToken {
	return LToken{
		Kind:  kind,
		Value: value,
	}
}

type Lexer struct {
	position int64
	text     string
	current  rune
}

func (l *Lexer) New(text string) *Lexer {
	s := &Lexer{
		text:     text,
		position: 0,
	}
	s.current = rune(s.text[uint64(s.position)])
	return s
}

func (l *Lexer) PutText(text string) {
	l.text = text
}

func (l *Lexer) advance() {
	l.position++
	if l.position >= int64(len(l.text)) {
		l.current = 0
	} else {
		l.current = rune(l.text[uint64(l.position)])
	}
}

func (l *Lexer) IsEOF() bool {
	return l.current == 0
}

func (l *Lexer) Next() *LToken {
	numberic := "0123456789"
	alphbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$_"
	brackets := "(){}[]<>"
	length := len(l.text)
	var token *LToken
	for idx := 0; idx <= length && l.current != 0; idx++ {
		if strings.ContainsRune(numberic, l.current) {
			result := ""
			for nidx := l.position; nidx < int64(len(l.text)) && strings.ContainsRune(numberic+".", l.current); nidx++ {
				result += string(l.current)
				l.advance()
			}
			token = &LToken{
				Kind:  NUMBER,
				Value: result,
			}
		} else if strings.ContainsRune(brackets, l.current) {
			switch l.current {
			case '<':
				token = &LToken{
					Kind:  LANGLE,
					Value: string(l.current),
				}
			case '>':
				token = &LToken{
					Kind:  RANGLE,
					Value: string(l.current),
				}
			case '{':
				token = &LToken{
					Kind:  LBRACE,
					Value: string(l.current),
				}
			case '}':
				token = &LToken{
					Kind:  RBRACE,
					Value: string(l.current),
				}
			case '[':
				token = &LToken{
					Kind:  LBRACKET,
					Value: string(l.current),
				}
			case ']':
				token = &LToken{
					Kind:  RBRACKET,
					Value: string(l.current),
				}
			}
		} else if strings.ContainsRune(alphbet, l.current) {
			result := ""
			for aidx := l.position; aidx < int64(len(l.text)) && strings.ContainsRune(alphbet+numberic, l.current); aidx++ {
				result += string(l.current)
				l.advance()
			}
			switch result {
			case "from":
				token = &LToken{
					Kind:  FROM,
					Value: result,
				}
			case "import":
				token = &LToken{
					Kind:  IMPORT,
					Value: result,
				}
			case "str":
				token = &LToken{
					Kind:  STR,
					Value: result,
				}
			case "enum":
				token = &LToken{
					Kind:  ENUM,
					Value: result,
				}
			case "type":
				token = &LToken{
					Kind:  TYPE,
					Value: result,
				}
			case "query":
				token = &LToken{
					Kind:  QUERY,
					Value: result,
				}
			case "body":
				token = &LToken{
					Kind:  BODY,
					Value: result,
				}
			case "bool":
				token = &LToken{
					Kind:  BOOL,
					Value: result,
				}
			case "get":
				token = &LToken{
					Kind:  GET,
					Value: result,
				}
			case "post":
				token = &LToken{
					Kind:  POST,
					Value: result,
				}
			case "option":
				token = &LToken{
					Kind:  OPTION,
					Value: result,
				}
			case "put":
				token = &LToken{
					Kind:  PUT,
					Value: result,
				}
			case "delete":
				token = &LToken{
					Kind:  DELETE,
					Value: result,
				}
			case "head":
				token = &LToken{
					Kind:  HEAD,
					Value: result,
				}
			case "patch":
				token = &LToken{
					Kind:  PATCH,
					Value: result,
				}
			case "required":
				token = &LToken{
					Kind: REQUIRED,
					Value:result,
				}
			case "optional":
				token = &LToken{
					Kind: OPTIONAL,
					Value: result,
				}
			case "deprecate":
				token = &LToken{
					Kind: DEPRECATE,
					Value: result,
				}
			default:
				token = &LToken{
					Kind:  IDENTIFIER,
					Value: result,
				}
			}
		} else if strings.ContainsRune("=", l.current) {
			token = &LToken{
				Kind:  ASSIGNMENT,
				Value: string(l.current),
			}
		} else if strings.ContainsRune(":", l.current) {
			token = &LToken{
				Kind:  COLON,
				Value: string(l.current),
			}
		} else if strings.ContainsRune(" \t", l.current) {
			// 跳过空白字符
			l.advance()
			continue
		} else if strings.ContainsRune("\"", l.current) {
			// 字符串
			result := ""
			for sidx := l.position; sidx < int64(len(l.text)) && l.current != '"'; sidx++ {
				result += string(l.current)
				l.advance()
			}
			token = &LToken{
				Kind:  STRING,
				Value: result,
			}
		} else {
			if l.current == '#' {
				// 忽略注释
				for cidx := l.position; cidx < int64(len(l.text)) && l.current != '\n'; cidx++ {
					l.advance()
					continue
				}
			} else {
				token = &LToken{
					Kind:  EOF,
					Value: "",
				}
			}
		}
		l.advance()
		return token
	}
	return &LToken{
		Kind:  EOF,
		Value: "",
	}
}
