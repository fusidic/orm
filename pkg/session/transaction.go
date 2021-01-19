package session

import "github.com/fusidic/orm/pkg/log"

// Begin a transcation.
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	// 调用 s.db.Begin() 得到 *sql.Tx 对象并赋值给 s.tx
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Commit a transaction.
func (s *Session) Commit() (err error) {
	log.Info("transcation commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

// Rollback a transaction.
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
