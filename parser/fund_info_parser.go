package parser

import (
	"Fachoi_fund_test/db_model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FundInfoParser struct {
}

func NewFundInfoParser() *FundInfoParser {
	return &FundInfoParser{}
}

func (flp *FundInfoParser) Parse(resp *http.Response) db_model.FundInfoModel {
	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fim := db_model.FundInfoModel{}

	dom.Find("table.info td").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" && selection.Text() != "---" {
			switch i {
			case 0:
				fim.FullName = selection.Text()
			case 1:
				fim.ShortName = selection.Text()
			case 2:
				reg := regexp.MustCompile(`\d+`)
				numbers := reg.FindAllString(selection.Text(), -1)
				fim.CodeFront = numbers[0]
				if len(numbers) == 2 {
					fim.CodeBack = numbers[1]
				}
			case 3:
				fim.FundType = selection.Text()
			case 4:
				// "2016年06月13日"
				tm, _ := time.Parse("2006年01月02日", selection.Text())
				fim.IssueDate = tm.Unix()
				fim.IssueDateString = selection.Text()
			case 5:
				// "2016年06月24日 / 29.380亿份"
				text := strings.Split(selection.Text(), " / ")[0]
				if text != "" {
					tm, _ := time.Parse("2006年01月02日", text)
					fim.LaunchDate = tm.Unix()
					fim.LaunchDateString = text
				}
			case 6:
				// "0.07亿元（截止至：2020年06月30日）"
				reg := regexp.MustCompile(`\d+\.?\d*`)
				numbers := reg.FindAllString(selection.Text(), -1)
				num1, _ := strconv.ParseFloat(numbers[0], 32)
				fim.AssetSize = float32(num1)
			case 8:
				fim.Company = selection.Text()
			case 9:
				fim.Trustee = selection.Text()
			case 10:
				fim.Manager = selection.Text()
			case 11:
				// "每份累计0.00元（0次）"
				reg := regexp.MustCompile(`\d+\.?\d*`)
				numbers := reg.FindAllString(selection.Text(), -1)
				num1, _ := strconv.ParseFloat(numbers[0], 32)
				num2, _ := strconv.ParseInt(numbers[1], 10, 0)
				fim.DividendPaymentPerUnit = float32(num1)
				fim.DividendCount = int(num2)
			default:
				// pass
			}
		}
	})

	tradeState := dom.Find(".col-right p:nth-child(2)").Text()
	// 个别页面无交易状态数据，如005135
	if tradeState != "" {
		tradeState = strings.ReplaceAll(tradeState, " ", "")
		tradeState = strings.ReplaceAll(tradeState, "\n", "")
		tradeState = strings.Split(tradeState, "：")[1]
	}
	fim.TradeState = tradeState
	return fim
}
