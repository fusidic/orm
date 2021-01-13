package session

import (
	"database/sql"
	"strings"

	"github.com/fusidic/orm/pkg/clause"
	"github.com/fusidic/orm/pkg/dialect"
	"github.com/fusidic/orm/pkg/log"
	"github.com/fusidic/orm/pkg/schema"
)

// Session is the structure to operate database.
type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

// New returns a session.
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// Clear reset sql Vars
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB ...
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw convert string to SQL
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

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db.
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
