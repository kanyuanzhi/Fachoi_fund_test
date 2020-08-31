package spider

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/saver"
	"database/sql"
	"fmt"
	"time"
)

type FundListSpider struct {
	*Spider
	parser *parser.FundListParser
	saver  *saver.FundListSaver
}

func NewFundListSpider(db *sql.DB) *FundListSpider {
	fls := new(FundListSpider)
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	fls.saver = saver.NewFundListSaver(db)
	return fls
}

func (fls *FundListSpider) Run() {
	dataChan := make(chan []db_model.FundListModel, 10)
	for {
		url, ok := fls.scheduler.Pop()
		if ok == false && len(dataChan) == 1 {
			fmt.Println("基金列表爬取完毕！等待存储完毕......")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		fls.process(url, dataChan)
	}
	fls.saver.Save(<-dataChan)
	close(dataChan)
	fmt.Println("基金列表存储完毕!")

}

func (fls *FundListSpider) process(url string, dataChan chan []db_model.FundListModel) {
	resp := fls.crawler.Crawl(url)
	if resp == nil {
		if !fls.crawler.UrlVisited(url) {
			fls.scheduler.Push(url)
		}
		return
	}
	dataChan <- fls.parser.Parse(resp)
}

// 获取所有前端基金代号（后端基金跳过）
func (fls *FundListSpider) GetFrontFundCodes() []string {
	return fls.parser.GetFrontFundCodes()
}
