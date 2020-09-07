package saver

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/util"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FundHistorySaver struct {
	DB *sqlx.DB
}

func NewFundHistorySaver(db *sqlx.DB) *FundHistorySaver {
	return &FundHistorySaver{
		DB: db,
	}
}

func (fhs *FundHistorySaver) Save(fhms []db_model.FundHistoryModel, code string) {
	util.CreateFundHistoryTable(code, fhs.DB)
	//util.TruncateTable("history_"+code+"_table", fhs.DB)
	if len(fhms) == 0 {
		// 无更新数据
		return
	} else if fhms[0].Date < 0 {
		// date字段出现负值，无历史数据
		return
	}

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
	_, err := fhs.DB.Exec(sqlStr, valueArgs...)
	util.CheckError(err, "FundHistorySaver")
}
