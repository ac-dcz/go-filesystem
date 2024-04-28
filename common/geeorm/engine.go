package geeorm

import (
	"database/sql"
	"log"
)

type Engine struct {
	pool *sql.DB //DBPools
}

func NewEngine(driver, dsn string) (*Engine, error) {
	pool, err := sql.Open(driver, dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Engine{pool}, nil
}

func (e *Engine) DB() *sql.DB {
	return e.pool
}

func (e *Engine) NewSession() *Session {
	return newSession(e.pool)
}
