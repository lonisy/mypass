package app

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

const (
	DataSourceName = "mypass.sqlite"
)

type DBStruct struct {
	db   *sql.DB
	once sync.Once
}

var Sqlite DBStruct

func (s *DBStruct) DB() *sql.DB {
	s.once.Do(func() {
		var err error
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Could not find local user folder. Error: %v\n", err)
		}
		s.db, err = sql.Open("sqlite3", userHomeDir+string(os.PathSeparator)+DataSourceName)
		if err != nil {
			panic(err.Error())
		}
	})
	return s.db
}
