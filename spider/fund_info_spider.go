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

func NewFundInfoSpider(threadsNum uint8, db *sql.DB) *FundInfoSpider {
	fis := new(FundInfoSpider)
	fis.Spider = NewSpider(threadsNum)
	fis.parser = parser.NewFundInfoParser()
	fis.saver = saver.NewFundInfoSaver(db)
	fis.Count = 1
	return fis
}

func (fis *FundInfoSpider) Run() {
	rm := NewResourceManagerChan(fis.threadsNum)
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
	go fis.saver.Save(fim)
}
