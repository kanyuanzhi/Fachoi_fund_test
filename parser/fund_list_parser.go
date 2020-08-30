package parser

import (
	"Fachoi_fund_test2/db_model"
	"io/ioutil"
	"net/http"
	"strings"
)

type FundListParser struct {
	fundCodes []string
}

func NewFundListParser() *FundListParser {
	return &FundListParser{
		fundCodes: make([]string, 0),
	}
}

func (flp *FundListParser) Parse(resp *http.Response) []db_model.FundListModel {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyBytes = bodyBytes[3:] // 去掉BOM头：0xef,0xbb,0xbf
	bodyStr := string(bodyBytes)
	bodyStr = strings.ReplaceAll(bodyStr, "var r = [[", "")
	bodyStr = strings.ReplaceAll(bodyStr, "]];", "")
	bodyStr = strings.ReplaceAll(bodyStr, "\"", "")

	var middleResult []string = strings.Split(bodyStr, "],[")
	var fundListData = make([]db_model.FundListModel, len(middleResult))

	for i, mr := range middleResult {
		item := strings.Split(mr, ",")
		flm := db_model.FundListModel{}
		flm.Code = item[0]
		flm.ShortName = item[2]
		flm.FundType = item[3]
		fundListData[i] = flm
		// 后端基金跳过
		if !strings.Contains(flm.ShortName, "后端") {
			flp.fundCodes = append(flp.fundCodes, flm.Code)
		}

	}
	return fundListData
}

// 获取所有前端基金代号（后端基金跳过）
func (flp *FundListParser) GetFundCodes() []string {
	return flp.fundCodes
}
