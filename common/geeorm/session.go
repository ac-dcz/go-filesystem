package geeorm

import (
	"database/sql"
	"log"
	"strings"
)

type Session struct {
	db        *sql.DB
	sql       strings.Builder
	vars      []any
	tablename string
}

func newSession(db *sql.DB) *Session {
	return &Session{
		db:   db,
		sql:  strings.Builder{},
		vars: make([]any, 0),
	}
}

func (s *Session) SetTableName(tname string) *Session {
	s.tablename = tname
	return s
}

func (s *Session) TableName() string {
	return s.tablename
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.vars = nil
}

func (s *Session) Raw(query string, argv ...any) *Session {
	s.sql.WriteString(query)
	s.sql.WriteByte(' ')
	s.vars = append(s.vars, argv...)
	return s
}

func (s *Session) Exex() (sql.Result, error) {
	defer s.Clear()
	log.Println(s.sql.String(), s.vars)
	res, err := s.db.Exec(s.sql.String(), s.vars...)
	if err != nil {
		log.Println(err)
	}
	return res, err
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Println(s.sql.String(), s.vars)
	return s.db.QueryRow(s.sql.String(), s.vars...)
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Println(s.sql.String(), s.vars)
	return s.db.Query(s.sql.String(), s.vars...)
}
