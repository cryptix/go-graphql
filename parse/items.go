package parse

import "fmt"

//go:generate stringer -type itemType
type itemType int

const (
	itemError itemType = iota

	itemObjName    // eg 'node'
	itemDot        // .
	itemFunction   // first
	itemLeftBrace  // (
	itemFnArgument // 123
	itemRightBrace // )
	itemLeftCurly  // {
	itemFieldName  // id
	itemComma      // ,
	itemRightCurly // }
	itemEOF
)

type item struct {
	typ itemType // Type,
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("<%-14s> %.10q...", i.typ, i.val)
	}
	return fmt.Sprintf("<%-14s> %q", i.typ, i.val)
}
