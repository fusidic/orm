package session

import (
	"reflect"

	"github.com/fusidic/orm/pkg/clause"
)

// Insert one or more records in database.
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).GetRefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 将对象 value 转换，并添加到 VALUES 中
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Find gets all eligible records and put them into objects.
func (s *Session) Find(values interface{}) error {
	// destSlice.Type().Elem() 获取切片的单个元素的类型 destType，
	// 使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，
	// 映射出表结构 RefTable()
	destSlice := reflect.Indirect((reflect.ValueOf(values)))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).GetRefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
