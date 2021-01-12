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
	log.Info("Connect datavase success")
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
