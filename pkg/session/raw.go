package session

import (
	"database/sql"
	"strings"

	"github.com/fusidic/orm/pkg/log"
)

// Session is the structure to operate database.
type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
}

// New returns a session.
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// Clear reset sql Vars
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// DB ...
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw ...
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
