package dbgenerate

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func (sql *sqlCli) getTables() ([]*tableInfo, error) {
	rows, err := SqlSelect(sql.db.DB(), "show tables")
	if err != nil {
		return nil, err
	}

	tbs := make([]*tableInfo, len(rows))
	for i, _ := range rows {
		tb := &tableInfo{}
		for _, v := range rows[i] {
			tb.Name = v
		}

		tb.Fields, err = sql.getTableInfo(tb.Name)
		if err != nil {
			fmt.Println("table", tb.Name, "err:", err)
		}
		tbs[i] = tb
	}

	return tbs, nil
}

func (sql *sqlCli) getTableInfo(tableName string) ([]*tableFieldInfo, error) {
	rows, err1 := sql.db.DB().Query(fmt.Sprintf("SHOW FULL FIELDS FROM %s", tableName))
	if err1 != nil {
		fmt.Println(err1.Error())
		return nil, err1
	}
	defer rows.Close()

	var tis []*tableFieldInfo
	for rows.Next() {
		var ti tableFieldInfo
		err2 := rows.Scan(&ti.Field, &ti.Type, &ti.Collation, &ti.Null, &ti.Key, &ti.Default, &ti.Extra, &ti.Privileges, &ti.Comment)
		if err2 != nil {
			fmt.Println(err2.Error())
			return nil, err2
		}
		tis = append(tis, &ti)
	}

	return tis, nil
}

func SqlSelect(db *sql.DB, sqlstr string, args ...interface{}) ([]map[string]string, error) {
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	var rets = make([]map[string]string, 0)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var ret = make(map[string]string)
		var value string
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}

		rets = append(rets, ret)
	}

	return rets, err
}

func newSqlCli(addr string) (*sqlCli, error) {
	db, err := gorm.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	return &sqlCli{db: db}, nil
}

type sqlCli struct {
	db *gorm.DB
}

type tableInfo struct {
	Name   string
	Fields []*tableFieldInfo
}

type tableFieldInfo struct {
	Field      string
	Type       string
	Collation  sql.NullString
	Null       string
	Key        string
	Default    sql.NullString
	Extra      string
	Privileges string
	Comment    string
}
