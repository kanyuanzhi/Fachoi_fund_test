package main

import (
	"Fachoi_fund_test2/db_mysql"
	"Fachoi_fund_test2/spider"
	"Fachoi_fund_test2/timer"
	"Fachoi_fund_test2/util"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// 换机运行前为避免数据库冲突，请在mysql_config.yaml文件中重新命名数据库
func main() {
	mysqlDB := db_mysql.NewMysql()
	mysqlDB.InitDatabase()
	db := mysqlDB.GetDB()

	if !util.InitializedCheck(db) {
		// 还未初始化所有表
		initializeAllTables(db)
	}

	update()
	t := timer.NewTimer(10)
	t.AddDayJob(update, 19, 0, 0)
	t.AddDayJob(update, 20, 0, 0)
	t.AddDayJob(update, 21, 0, 0)
	t.Run()
}

// 初始化所有表，爬取当前时间为止的所有基金相关数据
func initializeAllTables(db *sqlx.DB) {

	// 初始化基金列表表
	url := "http://fund.eastmoney.com/js/fundcode_search.js"
	fls := spider.NewFundListSpider(db)
	fls.AddUrl(url)
	fls.Run()

	//初始化基金信息表
	frontFundCodes := fls.GetFrontFundCodes()

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

	//初始化基金历史数据表（运行较慢）
	fhs := spider.NewFundHistorySpider(db, 5) // 此处并发数不能过大，否则数据存储有异常（与mysql所在服务器性能有关）
	fhs.AddUrls(historyUrls)
	fhs.Run()
}

// 更新基金数据
func update() {
	mysqlDB := db_mysql.NewMysql()
	db := mysqlDB.GetDB()
	url := "http://fund.eastmoney.com/js/fundcode_search.js"
	fls := spider.NewFundListSpider(db)
	fls.AddUrl(url)
	fls.Update()

	frontFundCodes := fls.GetFrontFundCodes()
	updatedFrontFundCodes := fls.GetUpdatedFrontFundCodes()
	var infoUrl, historyUrl string
	var infoUrls, historyUrls []string

	for _, code := range updatedFrontFundCodes {
		infoUrl = fmt.Sprintf("http://fundf10.eastmoney.com/jbgk_%s.html", code)
		infoUrls = append(infoUrls, infoUrl)
	}
	fis := spider.NewFundInfoSpider(db, 20)
	fis.AddUrls(infoUrls)
	fis.Run()

	for _, code := range frontFundCodes {
		historyUrl = fmt.Sprintf("http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery18307693431530679145_1599054971353"+
			"&fundCode=%s&pageIndex=1&pageSize=100&startDate=&endDate=&_=1599054971372", code)
		historyUrls = append(historyUrls, historyUrl)
	}
	fhs := spider.NewFundHistorySpider(db, 20)
	fhs.AddUrls(historyUrls)
	fhs.Update()
}
