package spider

import (
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/saver"
	"database/sql"
	"fmt"
	"time"
)

type FundInfoSpider struct {
	*Spider
	parser *parser.FundInfoParser
	saver  *saver.FundInfoSaver
	Count  int
}

func NewFundInfoSpider(db *sql.DB, threadsNum uint8) *FundInfoSpider {
	fis := new(FundInfoSpider)
	fis.Spider = NewSpider(threadsNum)
	fis.parser = parser.NewFundInfoParser()
	fis.saver = saver.NewFundInfoSaver(db)
	fis.Count = 1
	return fis
}

func (fis *FundInfoSpider) Run() {
	rm := NewResourceManager(fis.threadsNum)
	for {
		url, ok := fis.scheduler.Pop()
		if ok == false && rm.Has() == 0 {
			fmt.Println("url爬取完毕")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		rm.GetOne()
		go func(url string) {
			defer rm.FreeOne()
			fis.process(url)
			fis.Count++
			fmt.Println("FundInfoSpider: ", fis.Count)
		}(url)
	}
}

func (fis *FundInfoSpider) process(url string) {

	resp := fis.crawler.Crawl(url)
	if resp == nil {
		if !fis.crawler.UrlVisited(url) {
			fis.scheduler.Push(url)
		}
		return
	}
	//fis.parser.Parse(resp)

	fim := fis.parser.Parse(resp)
	//todo:这里用并发会导致数据还没存完主线程就已经结束，需要解决
	go fis.saver.Save(fim)
}
