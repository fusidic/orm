package schema

import (
	"go/ast"
	"reflect"

	"github.com/fusidic/orm/pkg/dialect"
)

// Schema is intended to deal with the convertion
// between object & table:
// type User struct {
// 	Name string `orm:"PRIMARY KEY"`
// 	Age  int
// }
// to
// CREATE TABLE `User` (`Name` text PRIMATY KEY, `Age` integer)

// Field represents a column of database.
type Field struct {
	Name string
	Type string
	Tag  string // 约束条件
}

// Schema represents a table of database.
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField ...
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse converts any objects to Schema.
func Parse(object interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(object)).Type()
	schema := &Schema{
		Model:    object,
		Name:     modelType.Name(), // 获取结构体的名称作为表名
		fieldMap: make(map[string]*Field),
	}
	// 获取实例的字段的个数
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 依次将 Object 中的元素转化为 sqlite 中对应的字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues returns the values of object's member variables.
// 将目标对象的成员变量平铺，如：将 &User{Name: "Tom", Age: 18} 转换为 ("Tom": 18)
func (schema *Schema) RecordValues(object interface{}) []interface{} {
	objectValue := reflect.Indirect(reflect.ValueOf(object))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, objectValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
