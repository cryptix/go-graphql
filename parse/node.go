package parse

import (
	"errors"
	"sort"
	"strconv"
)

type Node struct {
	name   string
	arg    string
	fields map[string]*Node
}

func (n Node) ID() (int64, error) {
	if n.arg == "" {
		return -1, errors.New("no id")
	}
	return strconv.ParseInt(n.arg, 10, 64)
}

func (n Node) Name() string {
	return n.name
}

func (n Node) FieldNames() []string {
	names := make([]string, len(n.fields))
	i := 0
	for _, v := range n.fields {
		names[i] = v.Name()
		i++
	}
	sort.Strings(names)
	return names
}

func (n *Node) Field(name string) (*Node, bool) {
	v, ok := n.fields[name]
	return v, ok
}

func (n *Node) PlainFields() []string {
	names := make([]string, len(n.fields))
	i := 0
	for _, v := range n.fields {
		if v.fields != nil {
			continue
		}
		names[i] = v.Name()
		i++
	}
	names = names[:i]
	sort.Strings(names)
	return names
}
