package app

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type DBStruct struct {
	db   *sql.DB
	once sync.Once
}

var Sqlite DBStruct

func (s *DBStruct) DB() *sql.DB {
	s.once.Do(func() {
		var err error
		s.db, err = sql.Open("sqlite3", "./passwords.db")
		if err != nil {
			panic(err.Error())
		}
	})
	return s.db
}
