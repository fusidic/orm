package clause

import "strings"

// Clause contains SQL conditions
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Type is the type of Clause
type Type int

// Support types for Clause
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

// Set adds a sub clause of specific type.
// Set 根据 Type 调用对应的 generator，并声称该子句对应的 SQL 语句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	// 调用对应数据操作，并传入参数
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build 根据传入 Type 的顺序，构造出最终的 SQL 语句
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
