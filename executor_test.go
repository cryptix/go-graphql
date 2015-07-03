package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet_simple(t *testing.T) {
	store := peopleStore{
		data: map[int64]person{
			123: person{123, "Frank", 23, "Green"},
			666: person{-1, "Devil", 2015, "Red"},
		},
	}
	exe := NewExecutor()
	if err := exe.Register("people", &store); err != nil {
		t.Fatalf("exe.Register() Err: %q", err)
	}
	cases := []struct {
		qry, res string
		err      bool
	}{
		{
			qry: `unknown(1){id,name}`,
			res: "store not registerd\n",
			err: true,
		},
		{
			qry: `people(1){id,name}`,
			res: "not found\n",
			err: true,
		},
		{
			qry: `people(123){id,name}`,
			res: `{"id":123,"name":"Frank"}` + "\n",
		},
		{
			qry: `people(123){age,haircolor}`,
			res: `{"age":23,"haircolor":"Green"}` + "\n",
		},
		{
			qry: `people(666){name,age}`,
			res: `{"age":2015,"name":"Devil"}` + "\n",
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest("GET", "/", strings.NewReader(c.qry))
		if err != nil {
			t.Fatalf("NewRequest Err: %q", err)
		}
		rw := httptest.NewRecorder()
		exe.ServeHTTP(rw, req)
		if !c.err {
			assert.Equal(t, http.StatusOK, rw.Code)
		}
		require.Equal(t, c.res, rw.Body.String())
	}
}

// totally naive test store
type person struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	HairColor string `json:"haircolor"`
}

type peopleStore struct {
	data map[int64]person
}

func (ps *peopleStore) Get(id int64, fields []string) (Data, error) {
	p, ok := ps.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	// lazy way to make map with only the required keys
	jsbytes, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}
	var tmp Data
	err = json.Unmarshal(jsbytes, &tmp)
	if err != nil {
		return nil, err
	}
	d := make(Data)
	for _, f := range fields {
		v, ok := tmp[f]
		if !ok {
			return nil, fmt.Errorf("graphql: Get required field not found %s", f)
		}
		d[f] = v
	}
	return d, nil
}
