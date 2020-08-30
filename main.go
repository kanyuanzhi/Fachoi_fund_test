package main

import (
	"Fachoi_fund_test2/db_mysql"
	"Fachoi_fund_test2/spider"
	"fmt"
	"strings"
)

func main() {
	mysqlDB := db_mysql.NewMysql()
	mysqlDB.InitDatabase()
	db := mysqlDB.GetDB()

	url := "http://fund.eastmoney.com/js/fundcode_search.js"
	fls := spider.NewFundListSpider()
	fls.AddUrl(url)
	fls.Run()

	fundCodes := fls.GetFundCodes()

	var urls []string
	for _, code := range fundCodes {
		code = strings.ReplaceAll(code, " ", "")
		url = fmt.Sprintf("http://fundf10.eastmoney.com/jbgk_%s.html", code)
		urls = append(urls, url)
	}
	fis := spider.NewFundInfoSpider(20, db)
	fis.AddUrls(urls)
	fis.Run()
}
