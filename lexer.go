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
	alphbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
					Kind:  TYPE,
					Value: result,
				}
			case "body":
				token = &LToken{
					Kind:  TYPE,
					Value: result,
				}
			case "bool":
				token = &LToken{
					Kind:  BOOL,
					Value: result,
				}
			case "GET":
				token = &LToken{
					Kind:  GET,
					Value: result,
				}
			case "POST":
				token = &LToken{
					Kind:  POST,
					Value: result,
				}
			case "OPTION":
				token = &LToken{
					Kind:  OPTION,
					Value: result,
				}
			case "PUT":
				token = &LToken{
					Kind:  PUT,
					Value: result,
				}
			case "DELETE":
				token = &LToken{
					Kind:  DELETE,
					Value: result,
				}
			case "HEAD":
				token = &LToken{
					Kind:  HEAD,
					Value: result,
				}
			case "PATCH":
				token = &LToken{
					Kind:  PATCH,
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
			l.advance()
			continue
		} else if strings.ContainsRune("\"", l.current) {
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
			token = &LToken{
				Kind:  EOF,
				Value: "",
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
