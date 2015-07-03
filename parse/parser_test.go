package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery_simple(t *testing.T) {
	n, err := Parse(`node(){id,name}`)
	require.Nil(t, err)
	require.NotNil(t, n)
	require.Equal(t, "node", n.Name())
	require.Len(t, n.FieldNames(), 2)
	require.Equal(t, []string{"id", "name"}, n.FieldNames())
}

func TestQuery_nested(t *testing.T) {
	r := require.New(t)
	n, err := Parse(`node(){id,name,obj{a,b}}`)
	r.Nil(err)
	r.NotNil(n)
	r.Equal("node", n.Name())
	r.Equal([]string{"id", "name", "obj"}, n.FieldNames())
	r.Equal([]string{"id", "name"}, n.PlainFields())
	obj, ok := n.Field("obj")
	r.True(ok)
	r.Equal([]string{"a", "b"}, obj.FieldNames())
	obj, ok = n.Field("nonExistant")
	r.False(ok)
	r.Nil(obj)
}
