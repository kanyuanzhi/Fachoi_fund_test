package db_model

type FundListModel struct {
	Id        int    `db:"id"`
	Code      string `db:"fund_code"`
	ShortName string `db:"fund_short_name"`
	FundType  string `db:"fund_type"`
}
