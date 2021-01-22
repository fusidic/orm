package orm

import (
	"database/sql"

	"github.com/fusidic/orm/pkg/dialect"
	"github.com/fusidic/orm/pkg/log"
	"github.com/fusidic/orm/pkg/session"
)

// Engine is the entrance of user
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine return a Engine
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return e, nil
}

// Close ...
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession encapsule session.New, returns a session.
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

// TxFunc is the interface for transaction
// the function imported will be called between tx.Begin() and tx.Commit()
// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
type TxFunc func(*session.Session) (interface{}, error)

// Transaction executes sql wrapped in a transaction, then automatically commit if no error occurs
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			rollbackErr := s.Rollback() // err is non-nil; don't change it
			log.Info(rollbackErr)
		} else {
			defer func() {
				if err != nil {
					rollbackErr := s.Rollback()
					log.Info(rollbackErr)
				}
			}()
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
