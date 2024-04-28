package store

import (
	"testing"
)

type Header struct {
	Name string
	Age  int
}

func (h Header) Hash() string {
	return h.Name
}

var datas = []Header{
	{"dcz", 23},
	{"sg", 22},
}

func TestSaveMemoryDB(t *testing.T) {
	db := NewMemoryDB[Header]("./meta/test.meta")
	if err := db.loadSource(); err != nil {
		t.Fatal(err)
	}
	for _, data := range datas {
		db.Put(data.Hash(), data)
	}
	for _, data := range datas {
		if val, ok := db.Get(data.Hash()); !ok {
			t.Fatalf("not found key")
		} else if val != data {
			t.Fatalf("orgin-data [%#v] data[%#v]", data, val)
		}
	}
}

func TestReadMemoryDB(t *testing.T) {
	db := NewMemoryDB[Header]("./meta/test.meta")
	if err := db.loadSource(); err != nil {
		t.Fatal(err)
	}
	for _, data := range datas {
		if val, ok := db.Get(data.Hash()); !ok {
			t.Fatalf("not found key")
		} else if val != data {
			t.Fatalf("orgin-data [%#v] data[%#v]", data, val)
		}
	}
}
