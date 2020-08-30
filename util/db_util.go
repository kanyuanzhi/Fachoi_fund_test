package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func TruncateTable(tableName string, db *sql.DB) {
	sqlStr := "TRUNCATE TABLE " + tableName
	_, err := db.Exec(sqlStr)
	CheckError(err, "TruncateTable")
}
