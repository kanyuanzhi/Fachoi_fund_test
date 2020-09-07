package db_mysql

import (
	"Fachoi_fund_test2/util"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type MysqlDB struct {
	db *sqlx.DB // 数据库连接池
}

func NewMysql() *MysqlDB {
	return &MysqlDB{}
}

// 初始化数据库，包括建立数据库，建立数据表
func (m *MysqlDB) InitDatabase() {
	createDatabase()
	m.db = createConnection()
	createTables(m.db)
}

// 获取连接
func (m *MysqlDB) GetDB() *sqlx.DB {
	if m.db == nil {
		fmt.Println("createConnection")
		m.db = createConnection()
	}
	return m.db
}

func createDatabase() {
	user, pass, host, port, _, _ := util.GetDBConfig()
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	_mysqlDb, _mysqlDbErr := sql.Open("mysql", dbDSN)
	if _mysqlDbErr != nil {
		panic("数据源配置不正确: " + _mysqlDbErr.Error())
	}

	sqlStr := "CREATE DATABASE IF NOT EXISTS fund_database CHARACTER SET utf8 COLLATE utf8_general_ci"

	_mysqlDb.Exec(sqlStr)
	_mysqlDb.Close()
}

func createConnection() *sqlx.DB {
	user, pass, host, port, dbname, charset := util.GetDBConfig()
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, pass, host, port, dbname, charset)
	db, err := sqlx.Connect("mysql", dbDSN)
	if err != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + err.Error())
	}
	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(120 * time.Second)

	if err = db.Ping(); nil != err {
		panic("数据库链接失败: " + err.Error())
	}
	return db
}

func createTables(db *sqlx.DB) {
	createFundListTable(db)
	createFundInfoTable(db)
}

func createFundListTable(db *sqlx.DB) {
	sqlStr := "CREATE TABLE IF NOT EXISTS fund_list_table (" +
		"id INT AUTO_INCREMENT, " +
		"fund_code VARCHAR(10), " +
		"fund_short_name VARCHAR(50), " +
		"fund_type VARCHAR(50), " +
		"PRIMARY KEY (id))"
	_, err := db.Exec(sqlStr)
	util.CheckError(err, "createFundListTable")
}

func createFundInfoTable(db *sqlx.DB) {
	sqlStr := "CREATE TABLE IF NOT EXISTS fund_info_table (" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"fund_code_front_end CHAR(6)," +
		"fund_code_back_end CHAR(6)," +
		"fund_full_name VARCHAR(50)," +
		"fund_short_name VARCHAR(50)," +
		"fund_type VARCHAR(50)," +
		"fund_issue_date BIGINT ," +
		"fund_issue_date_string VARCHAR(50)," +
		"fund_launch_date BIGINT ," +
		"fund_launch_date_string VARCHAR(50)," +
		"fund_asset_size FLOAT," +
		"fund_company VARCHAR(50)," +
		"fund_trustee VARCHAR(50)," +
		"fund_manager VARCHAR(50)," +
		"fund_dividend_payment_per_unit FLOAT," +
		"fund_dividend_count INT," +
		"fund_trade_state VARCHAR(50)" +
		")"
	_, err := db.Exec(sqlStr)
	util.CheckError(err, "createFundInfoTable")
}
