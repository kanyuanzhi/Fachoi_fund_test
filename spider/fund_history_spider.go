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

type FundHistorySpider struct {
	*Spider
	parser   *parser.FundHistoryParser
	saver    *saver.FundHistorySaver
	dataChan chan db_model.FundHistoryModelAndCode
}

func NewFundHistorySpider(db *sql.DB, threadsNum int) *FundHistorySpider {
	fhs := new(FundHistorySpider)
	fhs.Spider = NewSpider(threadsNum)
	fhs.parser = parser.NewFundHistoryParser()
	fhs.saver = saver.NewFundHistorySaver(db)
	fhs.dataChan = make(chan db_model.FundHistoryModelAndCode, 1000)
	return fhs
}

func (fhs *FundHistorySpider) Run() {
	crm := resource_manager.NewCrawlerResourceManager(fhs.threadsNum)
	srm := resource_manager.NewSaverResourceManager(fhs.urlsNum)
	srm.FillToTheFull()
	crawlFinished := false
	for {
		url, ok := fhs.scheduler.Pop()
		if crawlFinished == false && ok == false && crm.Has() == 0 {
			fmt.Println("基金历史数据爬取完毕！等待存储完毕......")
			crawlFinished = true
		} else if crawlFinished == false && ok == false {
			time.Sleep(time.Second)
			continue
		}
		crm.GetOne()
		// 并发爬取并解析页面
		go func(url string, crawlFinished bool) {
			defer crm.FreeOne()
			if crawlFinished {
				return
			}
			fhs.process(url)
			fhs.crawlCount <- 1
			fmt.Printf("爬取进度：%d / %d \n", len(fhs.crawlCount), fhs.urlsNum)
		}(url, crawlFinished)

		if crawlFinished && srm.Has() == 0 {
			fmt.Println("基金历史数据存储完毕！")
			break
		}
		// 并发存储数据
		go func() {
			defer srm.FreeOne()
			//fhs.saver.Save(<-fhs.dataChan)
			//fhs.saveCount <- 1
			//fmt.Printf("存储进度：%d / %d \n", len(fhs.saveCount), fhs.urlsNum)
		}()
	}

}

func (fhs *FundHistorySpider) process(url string) {
	resp := fhs.crawler.Crawl(url)
	if resp == nil {
		if !fhs.crawler.UrlVisited(url) {
			fhs.scheduler.Push(url)
		}
		<-fhs.crawlCount
		time.Sleep(time.Second * 5)
		return
	}
	// 识别出url中的基金代码
	//reg := regexp.MustCompile(`fundCode=\d+`)
	//code := strings.Split(reg.FindAllString(url, -1)[0], "=")[1]
	//fhs.dataChan <- db_model.FundHistoryModelAndCode{Fhms: fhs.parser.Parse(resp), Code: code}
}
