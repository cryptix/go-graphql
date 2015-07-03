/*
Package parse implements a lexer and parser for graphql - facebooks query language for relay/react

TODO

add functions of nodes - eg: .after()
*/
package parse

import (
	"errors"
	"fmt"
)

func Parse(qry string) (*Node, error) {
	n := new(Node)
	l := lex("base", qry)
parseLoop:
	for {
		i := <-l.items
		switch i.typ {
		case itemObjName:
			n = new(Node)
			n.fields = make(map[string]*Node)
			n.name = i.val
		case itemLeftBrace:
			if n == nil {
				return nil, errors.New("graphql: no root node")
			}
			i = <-l.items
			// no argument, return
			if i.typ == itemRightBrace {
				continue
			}
			if i.typ != itemFnArgument {
				return nil, errors.New("graphql: expected fnArgument")
			}
			n.arg = i.val
			i = <-l.items
			if i.typ != itemRightBrace {
				return nil, errors.New("graphql: expected rightBrace")
			}
		case itemLeftCurly:
			err := addNode(n, l.items)
			if err != nil {
				return nil, err
			}
		case itemError:
			return nil, fmt.Errorf("graphql: parse error: %s", i.val)
		case itemEOF:
			break parseLoop
		default:
			return nil, fmt.Errorf("graphql: parse - unhandled item type %v", i)
		}
	}
	return n, nil
}

func addNode(n *Node, items <-chan item) error {
	var newNode *Node
	for i := range items {
		switch i.typ {
		case itemComma:
		case itemFieldName:
			newNode = new(Node)
			newNode.name = i.val
			n.fields[i.val] = newNode
		case itemLeftCurly:
			if newNode == nil {
				return fmt.Errorf("graphql: addNode() - illegal leftCurly, field nil. %v", i)
			}
			newNode.fields = make(map[string]*Node)
			if err := addNode(newNode, items); err != nil {
				return err
			}
		case itemRightCurly:
			return nil
		default:
			return fmt.Errorf("graphql: addNode() - unhandled item type %v", i)
		}
	}
	panic("not reached")
}
