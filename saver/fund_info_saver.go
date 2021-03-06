package saver

import (
	"Fachoi_fund_test/db_model"
	"Fachoi_fund_test/util"
	"github.com/jmoiron/sqlx"
)

type FundInfoSaver struct {
	DB *sqlx.DB
}

func NewFundInfoSaver(db *sqlx.DB) *FundInfoSaver {
	return &FundInfoSaver{
		DB: db,
	}
}

func (fis *FundInfoSaver) Save(fim db_model.FundInfoModel) {
	sqlStr := "INSERT INTO fund_info_table(" +
		"fund_code_front_end, fund_code_back_end, fund_full_name, fund_short_name, fund_type, " +
		"fund_issue_date, fund_issue_date_string, fund_launch_date, fund_launch_date_string, " +
		"fund_asset_size, fund_company, fund_trustee, fund_manager, " +
		"fund_dividend_payment_per_unit, fund_dividend_count, fund_trade_state) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	smtp, err := fis.DB.Prepare(sqlStr)
	util.CheckError(err, "FundInfoSaver Save mysqlDB.Prepare")
	_, err = smtp.Exec(fim.CodeFront, fim.CodeBack, fim.FullName, fim.ShortName, fim.FundType,
		fim.IssueDate, fim.IssueDateString, fim.LaunchDate, fim.LaunchDateString,
		fim.AssetSize, fim.Company, fim.Trustee, fim.Manager,
		fim.DividendPaymentPerUnit, fim.DividendCount, fim.TradeState)
	util.CheckError(err, "FundInfoSaver Save smtp.Exec")
	smtp.Close()
}
