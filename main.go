package main

import (
	"Fachoi_fund_test2/db_mysql"
	"Fachoi_fund_test2/spider"
	"fmt"
)

func main() {
	mysqlDB := db_mysql.NewMysql()
	mysqlDB.InitDatabase()
	db := mysqlDB.GetDB()
	fmt.Println(*db)
	url := "http://fund.eastmoney.com/js/fundcode_search.js"
	fls := spider.NewFundListSpider(db)
	fls.AddUrl(url)
	fls.Run()

	frontFundCodes := fls.GetFrontFundCodes()
	fmt.Println(len(frontFundCodes))

	var infoUrls, historyUrls []string
	var infoUrl, historyUrl string
	for _, code := range frontFundCodes {
		infoUrl = fmt.Sprintf("http://fundf10.eastmoney.com/jbgk_%s.html", code)
		historyUrl = fmt.Sprintf("http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery18307693431530679145_1599054971353"+
			"&fundCode=%s&pageIndex=1&pageSize=10000&startDate=&endDate=&_=1599054971372", code)
		infoUrls = append(infoUrls, infoUrl)
		historyUrls = append(historyUrls, historyUrl)
	}
	fis := spider.NewFundInfoSpider(db, 20)
	fis.AddUrls(infoUrls)
	fis.Run()

	fhs := spider.NewFundHistorySpider(db, 20)
	fhs.AddUrls(historyUrls)
	fhs.Run()
}
