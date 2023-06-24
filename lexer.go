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

type Token struct {
	Kind  int
	Value string
}

func (c *Token) New(kind int, value string) Token {
	return Token{
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

func (l *Lexer) Next() *Token {
	numberic := "0123456789"
	alphbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$_"
	brackets := "(){}[]<>"
	length := len(l.text)
	var token *Token
	for idx := 0; idx <= length && l.current != 0; idx++ {
		if strings.ContainsRune(numberic, l.current) {
			result := ""
			for nidx := l.position; nidx < int64(len(l.text)) && strings.ContainsRune(numberic+".", l.current); nidx++ {
				result += string(l.current)
				l.advance()
			}
			token = &Token{
				Kind:  NUMBER,
				Value: result,
			}
		} else if strings.ContainsRune(brackets, l.current) {
			switch l.current {
			case '<':
				token = &Token{
					Kind:  LANGLE,
					Value: string(l.current),
				}
			case '>':
				token = &Token{
					Kind:  RANGLE,
					Value: string(l.current),
				}
			case '{':
				token = &Token{
					Kind:  LBRACE,
					Value: string(l.current),
				}
			case '}':
				token = &Token{
					Kind:  RBRACE,
					Value: string(l.current),
				}
			case '[':
				token = &Token{
					Kind:  LBRACKET,
					Value: string(l.current),
				}
			case ']':
				token = &Token{
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
				token = &Token{
					Kind:  FROM,
					Value: result,
				}
			case "import":
				token = &Token{
					Kind:  IMPORT,
					Value: result,
				}
			case "str":
				token = &Token{
					Kind:  STR,
					Value: result,
				}
			case "enum":
				token = &Token{
					Kind:  ENUM,
					Value: result,
				}
			case "type":
				token = &Token{
					Kind:  TYPE,
					Value: result,
				}
			case "query":
				token = &Token{
					Kind:  QUERY,
					Value: result,
				}
			case "body":
				token = &Token{
					Kind:  BODY,
					Value: result,
				}
			case "bool":
				token = &Token{
					Kind:  BOOL,
					Value: result,
				}
			case "GET":
				token = &Token{
					Kind:  GET,
					Value: result,
				}
			case "POST":
				token = &Token{
					Kind:  POST,
					Value: result,
				}
			case "OPTION":
				token = &Token{
					Kind:  OPTION,
					Value: result,
				}
			case "PUT":
				token = &Token{
					Kind:  PUT,
					Value: result,
				}
			case "DELETE":
				token = &Token{
					Kind:  DELETE,
					Value: result,
				}
			case "HEAD":
				token = &Token{
					Kind:  HEAD,
					Value: result,
				}
			case "PATCH":
				token = &Token{
					Kind:  PATCH,
					Value: result,
				}
			default:
				token = &Token{
					Kind:  IDENTIFIER,
					Value: result,
				}
			}
		} else if strings.ContainsRune("=", l.current) {
			token = &Token{
				Kind:  ASSIGNMENT,
				Value: string(l.current),
			}
		} else if strings.ContainsRune(":", l.current) {
			token = &Token{
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
			token = &Token{
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
				token = &Token{
					Kind:  EOF,
					Value: "",
				}
			}
		}
		l.advance()
		return token
	}
	return &Token{
		Kind:  EOF,
		Value: "",
	}
}
