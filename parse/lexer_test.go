package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex_empty(t *testing.T) {
	l := lex("empty", "")
	want := []item{
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

func TestLex_errors(t *testing.T) {
	tcases := []struct {
		name, qry string
		items     []item
	}{
		{"missing close brace",
			"node(",
			[]item{
				item{itemObjName, "node"},
				item{itemLeftBrace, "("},
				item{itemError, "illegal function argument"},
			},
		},
		// BUG(cryptix): error at parse level
		{"missing parens",
			"node{}",
			[]item{
				item{itemEOF, "node{}"},
				// item{itemError, "illegal parens argument"},
			},
		},
		{"missing close curly",
			"node(123){",
			[]item{
				item{itemObjName, "node"},
				item{itemLeftBrace, "("},
				item{itemFnArgument, "123"},
				item{itemRightBrace, ")"},
				item{itemLeftCurly, "{"},
				item{itemError, "illegal fieldname"},
			},
		},
	}
	for _, c := range tcases {
		l := lex(c.name, c.qry)
		assertLexer(t, l, c.items)
	}
}

func TestLex_simpleGraph(t *testing.T) {
	l := lex("simple", "node(123){one,two,obj{a,b}}")
	want := []item{
		item{itemObjName, "node"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "123"},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "one"},
		item{itemComma, ","},
		item{itemFieldName, "two"},
		item{itemComma, ","},
		item{itemFieldName, "obj"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "a"},
		item{itemComma, ","},
		item{itemFieldName, "b"},
		item{itemRightCurly, "}"},
		item{itemRightCurly, "}"},
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

func TestLex_indented(t *testing.T) {
	l := lex("indented", `node(123) {
		one,
		two,
		three
	}`)
	want := []item{
		item{itemObjName, "node"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "123"},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "one"},
		item{itemComma, ","},
		item{itemFieldName, "two"},
		item{itemComma, ","},
		item{itemFieldName, "three"},
		item{itemRightCurly, "}"},
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

func TestLex_talkExample(t *testing.T) {
	l := lex("talkExample",
		`node(1572451031) {
	id,
	name,
	birthdate {
		month,
		day
	},
	friends.after(3500401).first(2) {
		cursor,
		node {
			name
		}
	}
}`)
	want := []item{
		item{itemObjName, "node"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "1572451031"},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "id"},
		item{itemComma, ","},
		item{itemFieldName, "name"},
		item{itemComma, ","},
		item{itemFieldName, "birthdate"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "month"},
		item{itemComma, ","},
		item{itemFieldName, "day"},
		item{itemRightCurly, "}"},
		item{itemComma, ","},
		item{itemFieldName, "friends"},
		item{itemDot, "."},
		item{itemFunction, "after"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "3500401"},
		item{itemRightBrace, ")"},
		item{itemDot, "."},
		item{itemFunction, "first"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "2"},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "cursor"},
		item{itemComma, ","},
		item{itemFieldName, "node"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "name"},
		item{itemRightCurly, "}"},
		item{itemRightCurly, "}"},
		item{itemRightCurly, "}"},
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

// https://news.ycombinator.com/item?id=8978936
func TestLex_hnExample(t *testing.T) {
	l := lex("hnExample",
		`viewer() {
	posts {
		node {
			author { id, name, favorite_color },
		}
	},
	friends {
		node {
			id,
			name,
			favorite_color
		}
	},
	notifications {
		node {
			source { id, name, favorite_color }
		}
	},
}`)
	want := []item{
		item{itemObjName, "viewer"},
		item{itemLeftBrace, "("},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "posts"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "node"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "author"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "id"},
		item{itemComma, ","},
		item{itemFieldName, "name"},
		item{itemComma, ","},
		item{itemFieldName, "favorite_color"},
		item{itemRightCurly, "}"},
		item{itemComma, ","},
		item{itemRightCurly, "}"}, // node
		item{itemRightCurly, "}"}, // posts
		item{itemComma, ","},
		item{itemFieldName, "friends"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "node"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "id"},
		item{itemComma, ","},
		item{itemFieldName, "name"},
		item{itemComma, ","},
		item{itemFieldName, "favorite_color"},
		item{itemRightCurly, "}"}, // node
		item{itemRightCurly, "}"}, // friends
		item{itemComma, ","},
		item{itemFieldName, "notifications"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "node"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "source"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "id"},
		item{itemComma, ","},
		item{itemFieldName, "name"},
		item{itemComma, ","},
		item{itemFieldName, "favorite_color"},
		item{itemRightCurly, "}"}, // source
		item{itemRightCurly, "}"}, // node
		item{itemRightCurly, "}"}, // notifications
		item{itemComma, ","},
		item{itemRightCurly, "}"}, // viewer
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

func TestLex_ghcomment(t *testing.T) {
	l := lex("ghcomment", `friends.first(1) { edges { cursor, node { name } } }`)
	want := []item{
		item{itemObjName, "friends"},
		item{itemDot, "."},
		item{itemFunction, "first"},
		item{itemLeftBrace, "("},
		item{itemFnArgument, "1"},
		item{itemRightBrace, ")"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "edges"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "cursor"},
		item{itemComma, ","},
		item{itemFieldName, "node"},
		item{itemLeftCurly, "{"},
		item{itemFieldName, "name"},
		item{itemRightCurly, "}"},
		item{itemRightCurly, "}"},
		item{itemRightCurly, "}"},
		item{itemEOF, ""},
	}
	assertLexer(t, l, want)
}

func assertLexer(t *testing.T, l *lexer, want []item) {
	var got []item
	for i := range l.items {
		got = append(got, i)
	}
	require.Len(t, got, len(want), "delta: %d", len(got)-len(want))
	for idx := range want {
		assert.Equal(t,
			want[idx],
			got[idx],
			"item #%d from lexer is wrong\n Got:%s\nWant:%s", idx+1, got[idx], want[idx],
		)
	}
}
