package db_model

type FundInfoModel struct {
	Id                     int     `db:"id"`
	CodeFront              string  `db:"fund_code_front_end"`
	CodeBack               string  `db:"fund_code_back_end"`
	FullName               string  `db:"fund_full_name"`
	ShortName              string  `db:"fund_short_name"`
	FundType               string  `db:"fund_type"`
	IssueDate              int64   `db:"fund_issue_date"`
	IssueDateString        string  `db:"fund_issue_date_string"`
	LaunchDate             int64   `db:"fund_launch_date"`
	LaunchDateString       string  `db:"fund_launch_date_string"`
	AssetSize              float32 `db:"fund_asset_size"`
	Company                string  `db:"fund_company"`
	Trustee                string  `db:"fund_trustee"`
	Manager                string  `db:"fund_manager"`
	DividendPaymentPerUnit float32 `db:"fund_dividend_payment_per_unit"`
	DividendCount          int     `db:"fund_dividend_count"`
	TradeState             string  `db:"fund_trade_state"`
}
