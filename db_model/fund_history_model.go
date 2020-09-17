package db_model

type FundHistoryModel struct {
	Id               int     `db:id`
	Date             int64   `db:date`
	DateString       string  `db:date_string`
	Value            float32 `db:net_asset_value`
	AccumulatedValue float32 `db:accumulated_net_asset_value`
	Earnings         float32 `db:earnings_per_10000`
	AnnualReturn     float32 `db:7_day_annual_return`
}

type FundHistoryModelAndCode struct {
	Fhms []FundHistoryModel
	Code string
}

func NewFundHistoryModelAndCode(Fhms []FundHistoryModel, Code string) *FundHistoryModelAndCode {
	return &FundHistoryModelAndCode{
		Fhms: Fhms,
		Code: Code,
	}
}
