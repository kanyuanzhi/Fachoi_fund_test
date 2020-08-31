package spider

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/resource_manager"
	"Fachoi_fund_test2/saver"
	"database/sql"
	"fmt"
	"time"
)

type FundInfoSpider struct {
	*Spider
	parser *parser.FundInfoParser
	saver  *saver.FundInfoSaver
}

func NewFundInfoSpider(db *sql.DB, threadsNum int) *FundInfoSpider {
	fis := new(FundInfoSpider)
	fis.Spider = NewSpider(threadsNum)
	fis.parser = parser.NewFundInfoParser()
	fis.saver = saver.NewFundInfoSaver(db)
	return fis
}

func (fis *FundInfoSpider) Run() {
	crm := resource_manager.NewCrawlerResourceManager(fis.threadsNum)
	srm := resource_manager.NewSaverResourceManager(fis.urlsNum)
	srm.FillToTheFull()
	dataChan := make(chan db_model.FundInfoModel, 1000)
	crawlFinished := false
	for {
		url, ok := fis.scheduler.Pop()
		if crawlFinished == false && ok == false && crm.Has() == 0 {
			fmt.Println("基金信息爬取完毕！等待存储完毕......")
			crawlFinished = true
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		if crawlFinished && srm.Has() == 0 {
			fmt.Println("基金信息存储完毕！")
			break
		}
		crm.GetOne()

		// 并发爬取并解析页面
		go func(url string, dataChan chan db_model.FundInfoModel) {
			defer crm.FreeOne()
			fis.process(url, dataChan)
			fis.crawlCount <- 1
			fmt.Printf("爬取进度：%d / %d \n", len(fis.crawlCount), fis.urlsNum)
		}(url, dataChan)

		// 并发存储数据
		go func(dataChan chan db_model.FundInfoModel) {
			defer srm.FreeOne()
			fis.saver.Save(<-dataChan)
			fis.saveCount <- 1
			fmt.Printf("存储进度：%d / %d \n", len(fis.saveCount), fis.urlsNum)
		}(dataChan)
	}
}

func (fis *FundInfoSpider) process(url string, dataChan chan db_model.FundInfoModel) {
	resp := fis.crawler.Crawl(url)
	if resp == nil {
		if !fis.crawler.UrlVisited(url) {
			fis.scheduler.Push(url)
		}
		<-fis.crawlCount
		return
	}
	dataChan <- fis.parser.Parse(resp)
}
