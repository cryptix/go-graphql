package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	name  string // used only for err
	input string // the string beeing scanned
	state stateFn
	pos   int
	start int
	width int
	items chan item
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexObjName; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t itemType) {
	val := l.input[l.start:l.pos]
	// BUG(cryptix): a bit hacky...
	if t == itemFieldName || t == itemRightCurly {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			return
		}
	}
	l.items <- item{t, val}
	l.start = l.pos
}

// helpers
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

const (
	dot        = '.'
	comma      = ','
	leftBrace  = '('
	rightBrace = ')'
	leftCurly  = '{'
	rightCurly = '}'
)

func lexObjName(l *lexer) stateFn {
	for len(l.input[l.pos:]) > 0 {
		switch l.input[l.pos] {
		case leftBrace:
			if l.pos > l.start {
				l.emit(itemObjName) // emit identifier e.g. node
			}
			return lexLeftBrace
		case dot:
			if l.pos > l.start {
				l.emit(itemObjName) // emit identifier e.g. node
			}
			return lexDot
		}
		if l.next() == eof {
			break
		}
	}
	l.emit(itemEOF) // empty query
	return nil
}

func lexDot(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemDot)
	return lexFnName
}

func lexFnName(l *lexer) stateFn {
	for {
		if l.input[l.pos] == leftBrace {
			if l.pos > l.start {
				l.emit(itemFunction)
			}
			return lexLeftBrace
		}
		if l.next() == eof {
			return l.errorf("illegal function name")
		}
	}
	panic("not reached")
}

func lexLeftBrace(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemLeftBrace)
	return lexFnArgument
}

func lexFnArgument(l *lexer) stateFn {
	for len(l.input[l.pos:]) > 0 {
		if l.input[l.pos] == rightBrace {
			if l.pos > l.start {
				l.emit(itemFnArgument) // emit function argument as string
			}
			return lexRightBrace
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			break
		}
	}
	return l.errorf("illegal function argument")
}

func lexRightBrace(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemRightBrace)
	l.acceptRun(" \t\n")
	l.ignore()
	if l.peek() == dot {
		return lexDot
	}
	return lexLeftCurly
}

func lexLeftCurly(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemLeftCurly)
	p := l.peek()
	for isSpace(p) || p == '\n' {
		l.next()
		l.ignore()
		p = l.peek()
	}
	return lexFieldNames
}

func lexFieldNames(l *lexer) stateFn {
	for len(l.input[l.pos:]) > 0 {
		switch l.input[l.pos] {
		case rightCurly:
			if l.pos > l.start {
				l.emit(itemFieldName)
			}
			return lexRightCurly
		case leftCurly:
			if l.pos > l.start {
				l.emit(itemFieldName)
			}
			return lexLeftCurly
		case comma:
			if l.pos > l.start {
				l.emit(itemFieldName)
			}
			return lexComma
		case dot:
			if l.pos > l.start {
				l.emit(itemFieldName)
			}
			return lexDot
		case leftBrace:
			if l.pos > l.start {
				l.emit(itemFunction)
			}
			return lexLeftBrace
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("illegal fieldname")
}

func lexComma(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemComma)
	p := l.peek()
	for isSpace(p) || p == '\n' {
		l.next()
		l.ignore()
		p = l.peek()
	}
	return lexFieldNames
}

func lexRightCurly(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemRightCurly)
	// TODO add depth check
	if l.peek() == eof {
		l.emit(itemEOF)
		return nil
	}
	return lexFieldNames
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
