package orm

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fusidic/orm/pkg/session"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "../../orm.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func Test_NewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func Test_Engine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		testTransactionRollback(t)
	})
	// t.Run("commit", func(t *testing.T) {
	// 	testTransactionCommit(t)
	// })
}

func testTransactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func testTransactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, err := s.Insert(&User{"Tom", 18})
		return nil, err
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

// test Migrate
func Test_Engine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values(?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
