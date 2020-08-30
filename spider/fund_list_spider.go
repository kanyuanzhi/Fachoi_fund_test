package spider

import (
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/saver"
	"database/sql"
	"fmt"
	"time"
)

type FundListSpider struct {
	*Spider
	parser    *parser.FundListParser
	saver     *saver.FundListSaver
	fundCodes *[]string
}

func NewFundListSpider(db *sql.DB) *FundListSpider {
	fls := new(FundListSpider)
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	fls.saver = saver.NewFundListSaver(db)
	return fls
}

func (fls *FundListSpider) Run() {
	rm := NewResourceManager(fls.threadsNum)
	for {
		url, ok := fls.scheduler.Pop()
		if ok == false && rm.Has() == 0 {
			fmt.Println("爬取完毕!")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		rm.GetOne()
		go func(url string) {
			defer rm.FreeOne()
			fls.process(url)
		}(url)
	}
}

func (fls *FundListSpider) process(url string) {
	resp := fls.crawler.Crawl(url)
	if resp == nil {
		if !fls.crawler.UrlVisited(url) {
			fls.scheduler.Push(url)
		}
		return
	}
	flms := fls.parser.Parse(resp)
	fls.saver.Save(flms)
}

// 获取所有前端基金代号（后端基金跳过）
func (fls *FundListSpider) GetFundCodes() []string {
	return fls.parser.GetFundCodes()
}
