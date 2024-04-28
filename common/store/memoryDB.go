package store

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const DefaultMetaPath = "./common/store/meta/memory.meta"

type MemoryDB[T KeyValue] struct {
	dsn   string
	mu    sync.RWMutex
	datas map[string]T
}

func NewMemoryDB[T KeyValue](dsn string) *MemoryDB[T] {
	db := &MemoryDB[T]{
		mu:    sync.RWMutex{},
		datas: make(map[string]T),
		dsn:   dsn,
	}
	if err := db.loadSource(); err != nil {
		log.Println(err)
	}
	// runtime.SetFinalizer(db, func(db *MemoryDB[T]) {
	// 	db.SaveSource()
	// })
	return db
}

func (db *MemoryDB[T]) loadSource() error {
	now := time.Now()
	f, err := os.OpenFile(db.dsn, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	for {
		val := new(T)
		err := decoder.Decode(val)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		db.datas[(*val).Hash()] = *val
	}
	log.Printf("memory database init source finish [time %.3fms]\n", time.Since(now).Seconds())
	return nil
}

func (db *MemoryDB[T]) saveSource() error {
	// db.mu.RLock()
	// defer db.mu.RUnlock()
	f, err := os.OpenFile(db.dsn, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	for _, val := range db.datas {
		if err := encoder.Encode(val); err != nil {
			f.Truncate(int64(io.SeekCurrent))
			return err
		}
	}
	return nil
}

func (db *MemoryDB[T]) Get(key string) (T, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.datas[key]
	return val, ok
}

func (db *MemoryDB[T]) Put(key string, val T) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.datas[key] = val
	db.saveSource()
}

func (db *MemoryDB[T]) Keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, 0)
	for key := range db.datas {
		keys = append(keys, key)
	}
	return keys
}
