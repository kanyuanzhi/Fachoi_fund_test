package update_checker

import (
	"Fachoi_fund_test/db_model"
	"Fachoi_fund_test/util"
	"github.com/jmoiron/sqlx"
)

type FundHistoryUpdateChecker struct {
	db                    *sqlx.DB
	updatedFrontFundCodes []string
}

func NewFundHistoryUpdateChecker(db *sqlx.DB) *FundHistoryUpdateChecker {
	return &FundHistoryUpdateChecker{
		db:                    db,
		updatedFrontFundCodes: make([]string, 0),
	}
}

func (fhuc *FundHistoryUpdateChecker) Check(fhms []db_model.FundHistoryModel, code string) []db_model.FundHistoryModel {
	util.CreateFundHistoryTable(code, fhuc.db) //新基金没有对应历史数据表，需要先建立
	sqlStr := "select date from history_" + code + "_table order by id desc limit 1"
	rows, err := fhuc.db.Query(sqlStr)
	util.CheckError(err, "FundHistoryUpdateChecker Check select")
	var latestDate int64
	for rows.Next() {
		rows.Scan(&latestDate)
		if latestDate == 0 { //历史数据表为空表
			rows.Close()
			return fhms // 不做切片直接返回所有fhms，可能包含date字段为负值的数据，在saver中对其做处理
		}
	}
	rows.Close()

	for i, fhm := range fhms {
		if fhm.Date == latestDate {
			fhms = fhms[i+1:]
		}
	}
	return fhms // 可能为空，在saver中对其做处理
}
