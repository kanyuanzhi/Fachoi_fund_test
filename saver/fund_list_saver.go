package saver

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/util"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FundListSaver struct {
	db *sqlx.DB
}

func NewFundListSaver(db *sqlx.DB) *FundListSaver {
	util.TruncateTable("fund_list_table", db)
	return &FundListSaver{
		db: db,
	}
}

func (fls *FundListSaver) Save(flms []db_model.FundListModel) {
	flmsSize := len(flms)
	valueStrings := make([]string, 0, flmsSize)
	valueArgs := make([]interface{}, 0, 3*flmsSize)
	for _, flm := range flms {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, flm.Code)
		valueArgs = append(valueArgs, flm.ShortName)
		valueArgs = append(valueArgs, flm.FundType)
	}
	sqlStr := fmt.Sprintf("INSERT INTO fund_list_table(fund_code, fund_short_name, fund_type) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := fls.db.Exec(sqlStr, valueArgs...)
	util.CheckError(err, "FundListSaver")
}
