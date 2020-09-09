package parser

import (
	"Fachoi_fund_test/db_model"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FundHistoryParser struct {
}

func NewFundHistoryParser() *FundHistoryParser {
	return &FundHistoryParser{}
}

func (flp *FundHistoryParser) Parse(resp *http.Response) []db_model.FundHistoryModel {
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ReplaceAll(bodyStr, "jQuery183019601346852042933_1596811572354(", "")
	bodyStr = strings.ReplaceAll(bodyStr, ")", "")

	// 基金历史数据包括两种格式，一种为单位净值+累计净值，一种为每万份收益+7日年化收益率
	var isValueFlag bool
	if gjson.Get(bodyStr, "Data.SYType").Raw == "null" {
		isValueFlag = true
	} else {
		isValueFlag = false
	}
	historyData := gjson.Get(bodyStr, "Data.LSJZList")

	historyData.Raw = strings.ReplaceAll(historyData.Raw, "[", "")
	historyData.Raw = strings.ReplaceAll(historyData.Raw, "]", "")
	historyDataArray := strings.Split(historyData.Raw, "},{")
	var fhms []db_model.FundHistoryModel

	for _, item := range historyDataArray {
		// strings.Split(historyData.Raw, "},{")时祛除了原item中的{}符号，此处加上以供gjson读取（添加左侧{即可）
		item = "{" + item
		fhm := db_model.FundHistoryModel{}
		tm, _ := time.Parse("2006-01-02", gjson.Get(item, "FSRQ").Str)
		fhm.Date = tm.Unix()
		fhm.DateString = gjson.Get(item, "FSRQ").Str
		valueString := gjson.Get(item, "DWJZ").Str
		accValueString := gjson.Get(item, "LJJZ").Str
		if isValueFlag {
			if valueString == "" {
				fhm.Value = -1 // 历史数据中没有当天数据，设为-1
			} else {
				f, _ := strconv.ParseFloat(valueString, 32)
				fhm.Value = float32(f)
			}
			if accValueString == "" {
				fhm.AccumulatedValue = -1
			} else {
				f, _ := strconv.ParseFloat(accValueString, 32)
				fhm.AccumulatedValue = float32(f)
			}
			fhm.Earnings = 0
			fhm.AnnualReturn = 0
		} else {
			if valueString == "" {
				fhm.Earnings = -1
			} else {
				f, _ := strconv.ParseFloat(valueString, 32)
				fhm.Earnings = float32(f)
			}
			if accValueString == "" {
				fhm.AnnualReturn = -1
			} else {
				f, _ := strconv.ParseFloat(accValueString, 32)
				fhm.AnnualReturn = float32(f)
			}
			fhm.Value = 0
			fhm.AccumulatedValue = 0
		}
		var temp = []db_model.FundHistoryModel{fhm}
		fhms = append(temp, fhms...) // 在切片头部插入
	}
	return fhms
}
