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
	position uint64
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
	if l.position >= uint64(len(l.text)) {
		l.current = 0
	} else {
		l.current = rune(l.text[l.position])
	}
}

func (l *Lexer) IsEOF() bool {
	return l.current == 0
}

func (l *Lexer) Next() *Token {
	numberic := "0123456789"
	alphbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$_"
	brackets := "(){}[]<>"
	length := uint64(len(l.text))
	var token *Token
	for idx  := uint64(0); idx <= length && l.current != 0; idx++ {
		if strings.ContainsRune(numberic, l.current) {
			result := ""
			for nidx := l.position; nidx < length && strings.ContainsRune(numberic+".", l.current); nidx++ {
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
			for aidx := l.position; aidx < length && strings.ContainsRune(alphbet+numberic, l.current); aidx++ {
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
			case "get":
				token = &Token{
					Kind:  GET,
					Value: result,
				}
			case "post":
				token = &Token{
					Kind:  POST,
					Value: result,
				}
			case "option":
				token = &Token{
					Kind:  OPTION,
					Value: result,
				}
			case "put":
				token = &Token{
					Kind:  PUT,
					Value: result,
				}
			case "delete":
				token = &Token{
					Kind:  DELETE,
					Value: result,
				}
			case "head":
				token = &Token{
					Kind:  HEAD,
					Value: result,
				}
			case "patch":
				token = &Token{
					Kind:  PATCH,
					Value: result,
				}
			case "required":
				token = &Token{
					Kind: REQUIRED,
					Value:result,
				}
			case "optional":
				token = &Token{
					Kind: OPTIONAL,
					Value: result,
				}
			case "deprecate":
				token = &Token{
					Kind: DEPRECATE,
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
			for sidx := l.position; sidx < length && l.current != '"'; sidx++ {
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
				for cidx := l.position; cidx < length && l.current != '\n'; cidx++ {
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
