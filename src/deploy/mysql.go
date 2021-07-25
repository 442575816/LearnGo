package deploy

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBUtils struct {
	db *sql.DB
}

// Open 打开连接
func (dbutils *DBUtils) Open(db string, host string, port int, user string, password string) {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, db)

	dbInstance, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	dbutils.db = dbInstance
}

func (dbutils *DBUtils) Query(sql string) []map[string]string {
	rows, err := dbutils.db.Query(sql)
	defer rows.Close()

	if err != nil {
		fmt.Printf("query data error:%v\n", err)
		return nil
	}

	var results []map[string]string
	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		results = append(results, record)
	}

	return results
}

func (dbutils *DBUtils) Update(sql string) int64 {
	result, err := dbutils.db.Exec(sql)
	if err != nil {
		fmt.Println("update failed, error:", err.Error())
		return 0
	}

	rows, _ := result.RowsAffected()
	return rows
}

func (dbutils *DBUtils) Insert(sql string) int64 {
	result, err := dbutils.db.Exec(sql)
	if err != nil {
		fmt.Println("update failed, error:", err.Error())
		return 0
	}

	lastId, _ := result.LastInsertId()
	return lastId
}

func (dbutils *DBUtils) Delete(sql string) int64 {
	result, err := dbutils.db.Exec(sql)
	if err != nil {
		fmt.Println("update failed, error:", err.Error())
		return 0
	}

	rows, _ := result.RowsAffected()
	return rows
}

func (dbutils *DBUtils) Close() {
	if dbutils.db != nil {
		dbutils.db.Close()
	}
}
