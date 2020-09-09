package update_checker

import (
	"Fachoi_fund_test/db_model"
	"Fachoi_fund_test/util"
	"crypto/md5"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FundListUpdateChecker struct {
	db                    *sqlx.DB
	updatedFrontFundCodes []string
}

func NewFundListUpdateChecker(db *sqlx.DB) *FundListUpdateChecker {
	return &FundListUpdateChecker{
		db:                    db,
		updatedFrontFundCodes: make([]string, 0),
	}
}

func (fluc *FundListUpdateChecker) Check(flms []db_model.FundListModel) []db_model.FundListModel {
	sqlStr := "select fund_code from fund_list_table"
	rows, err := fluc.db.Query(sqlStr)
	util.CheckError(err, "FundListUpdateChecker Check")
	var code string
	OldCodesMap := make(map[[md5.Size]byte]string)
	for rows.Next() {
		rows.Scan(&code)
		key := md5.Sum([]byte(code))
		OldCodesMap[key] = code
	}
	rows.Close()

	var candidateFLMs []db_model.FundListModel
	for _, flm := range flms {
		key := md5.Sum([]byte(flm.Code))
		if _, has := OldCodesMap[key]; has == false {
			candidateFLMs = append(candidateFLMs, flm)
			if !strings.Contains(flm.ShortName, "后端") {
				fluc.updatedFrontFundCodes = append(fluc.updatedFrontFundCodes, flm.Code)
			}
		}
	}
	return candidateFLMs
}

func (fluc *FundListUpdateChecker) GetUpdatedFrontFundCodes() []string {
	return fluc.updatedFrontFundCodes
}
