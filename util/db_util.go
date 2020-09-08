package util

import (
	"Fachoi_fund_test2/resource_manager"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

func TruncateTable(tableName string, db *sqlx.DB) {
	sqlStr := "TRUNCATE TABLE " + tableName
	_, err := db.Exec(sqlStr)
	CheckError(err, "TruncateTable")
}

func CreateFundHistoryTable(code string, db *sqlx.DB) int64 {
	tableName := "history_" + code + "_table"
	sqlStr := "CREATE TABLE IF NOT EXISTS " + tableName + " (" +
		"id INT AUTO_INCREMENT," +
		"date BIGINT," +
		"date_string VARCHAR(50)," +
		"net_asset_value FLOAT, " +
		"accumulated_net_asset_value FLOAT, " +
		"earnings_per_10000 FLOAT, " +
		"7_day_annual_return FLOAT, " +
		"PRIMARY KEY (id)" +
		")"
	result, err := db.Exec(sqlStr)

	CheckError(err, "createFundInfoTable")
	rowNum, _ := result.RowsAffected()
	return rowNum
}

func CreateAllFundHistoryTables(db *sqlx.DB, codes []string) {
	createThreadsManager := resource_manager.NewResourceManager(20)
	for _, code := range codes {
		createThreadsManager.GetOne()
		go func() {
			defer createThreadsManager.FreeOne()
			CreateFundHistoryTable(code, db)
		}()
		if createThreadsManager.Has() == 0 {
			fmt.Println("所有基金历史数据表创建完毕!")
			break
		}
	}
}

//func CheckTableExist(tableName string, db *sql.DB) bool{
//	sqlStr := "TRUNCATE TABLE " + tableName
//	_, err := db.Exec(sqlStr)
//	CheckError(err, "TruncateTable")
//}

func EvictExceptionData(db *sqlx.DB) {
	sqlStr := "SELECT table_name FROM information_schema.TABLES"
	rows, _ := db.Query(sqlStr)
	var tableName string
	var tableNames []string

	for rows.Next() {
		rows.Scan(&tableName)
		if strings.Contains(tableName, "history") {
			tableNames = append(tableNames, tableName)
		}
	}
	rows.Close()

	erm := resource_manager.NewResourceManager(20)
	for _, tableName := range tableNames {
		erm.GetOne()
		go func(tableName string) {
			defer erm.FreeOne()
			var latestDate int64
			sqlStr = "select `date` from " + tableName + " order by id desc limit 1"
			rows2, _ := db.Query(sqlStr)
			for rows2.Next() {
				rows2.Scan(&latestDate)
				if latestDate < 0 {
					TruncateTable(tableName, db)
					fmt.Println(tableName)
				}
			}
			rows2.Close()
		}(tableName)
		if erm.Has() == 0 {
			return
		}
	}
}

func InitializedCheck(db *sqlx.DB) bool {
	sqlStr := "select id from fund_list_table limit 1"
	var id int
	rows, _ := db.Query(sqlStr)
	for rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	if id != 0 {
		return true
	} else {
		return false
	}
}
