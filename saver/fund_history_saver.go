package saver

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/db_mysql"
	"Fachoi_fund_test2/util"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type FundHistorySaver struct {
	db *sql.DB
}

func NewFundHistorySaver(db *sql.DB) *FundHistorySaver {
	return &FundHistorySaver{
		db: db,
	}
}

func (fhs *FundHistorySaver) Save(fhms []db_model.FundHistoryModel, code string) {
	//fmt.Println(fhs.db.Stats())
	db_mysql.CreateFundHistoryTable(fhs.db, code)
	//util.TruncateTable("history_"+code+"_table", fhs.db)
	fhmsSize := len(fhms)
	valueStrings := make([]string, 0, fhmsSize)
	valueArgs := make([]interface{}, 0, 6*fhmsSize)
	for _, fhm := range fhms {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, fhm.Date)
		valueArgs = append(valueArgs, fhm.DateString)
		valueArgs = append(valueArgs, fhm.Value)
		valueArgs = append(valueArgs, fhm.AccumulatedValue)
		valueArgs = append(valueArgs, fhm.Earnings)
		valueArgs = append(valueArgs, fhm.AnnualReturn)
	}

	sqlStr := fmt.Sprintf("INSERT INTO history_"+code+"_table"+
		"(date, date_string, net_asset_value, accumulated_net_asset_value,earnings_per_10000, 7_day_annual_return) "+
		"VALUES %s",
		strings.Join(valueStrings, ","))
	t1 := time.Now().Unix()
	_, err := fhs.db.Exec(sqlStr, valueArgs...)
	t2 := time.Now().Unix()
	fmt.Println(len(valueArgs))
	fmt.Println(t2 - t1)
	util.CheckError(err, "FundHistorySaver")
}
