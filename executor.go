package graphql

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cryptix/go-graphql/parse"
)

type Data map[string]interface{}

type Store interface {
	Get(int64, []string) (Data, error)
}

type Executor struct {
	stores map[string]Store
}

func NewExecutor() *Executor {
	return &Executor{
		stores: make(map[string]Store),
	}
}

func (e *Executor) Register(name string, s Store) error {
	// BUG(cryptix) dont replace
	e.stores[name] = s
	return nil
}

func (e *Executor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	txt, err := ioutil.ReadAll(io.LimitReader(req.Body, 512*1024)) // who would ever need more than half a meg of query string... ;)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	qry, err := parse.Parse(string(txt))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	store, ok := e.stores[qry.Name()]
	if !ok {
		http.Error(rw, "store not registerd", http.StatusNotFound)
		return
	}
	id, err := qry.ID()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := store.Get(id, qry.PlainFields())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
