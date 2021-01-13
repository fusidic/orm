package session

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fusidic/orm/pkg/log"
	"github.com/fusidic/orm/pkg/schema"
)

// Model 方法用于给 refTable 赋值
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// GetRefTable returns a Schema instance that contains all parsed fields.
func (s *Session) GetRefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable create a table in database with model.
func (s *Session) CreateTable() error {
	table := s.GetRefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable drop a table in database
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.GetRefTable().Name)).Exec()
	return err
}

// HasTable check if the database has the table
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.GetRefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.GetRefTable().Name
}
