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
	BOOL              // bool
	FROM              // from
	IMPORT            // import
	COLON             //:
	IDENTIFIER        // identifier
	ASSIGNMENT        // =
	EOF               // end
)

type LToken struct {
	kind  int
	value string
}

func (c *LToken) New(kind int, value string) LToken {
	return LToken{
		kind:  kind,
		value: value,
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
				kind:  NUMBER,
				value: result,
			}
		} else if strings.ContainsRune(brackets, l.current) {
			switch l.current {
			case '<':
				token = &LToken{
					kind:  LANGLE,
					value: string(l.current),
				}
			case '>':
				token = &LToken{
					kind:  RANGLE,
					value: string(l.current),
				}
			case '{':
				token = &LToken{
					kind:  LBRACE,
					value: string(l.current),
				}
			case '}':
				token = &LToken{
					kind:  RBRACE,
					value: string(l.current),
				}
			case '[':
				token = &LToken{
					kind:  LBRACKET,
					value: string(l.current),
				}
			case ']':
				token = &LToken{
					kind:  RBRACKET,
					value: string(l.current),
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
					kind:  FROM,
					value: result,
				}
			case "import":
				token = &LToken{
					kind:  IMPORT,
					value: result,
				}
			case "str":
				token = &LToken{
					kind:  STR,
					value: result,
				}
			case "enum":
				token = &LToken{
					kind:  ENUM,
					value: result,
				}
			case "type":
				token = &LToken{
					kind:  TYPE,
					value: result,
				}
			default:
				token = &LToken{
					kind:  IDENTIFIER,
					value: result,
				}
			}
		} else if strings.ContainsRune("=", l.current) {
			token = &LToken{
				kind:  ASSIGNMENT,
				value: string(l.current),
			}
		} else if strings.ContainsRune(":", l.current) {
			token = &LToken{
				kind:  COLON,
				value: string(l.current),
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
				kind:  STRING,
				value: result,
			}
		} else {
			token = &LToken{
				kind:  EOF,
				value: "",
			}
		}
		l.advance()
		return token
	}
	return &LToken{
		kind:  EOF,
		value: "",
	}
}
