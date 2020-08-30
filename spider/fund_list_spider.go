package spider

import (
	"Fachoi_fund_test2/parser"
	"fmt"
	"time"
)

type FundListSpider struct {
	*Spider
	parser    *parser.FundListParser
	fundCodes *[]string
}

func NewFundListSpider() *FundListSpider {
	fls := new(FundListSpider)
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	return fls
}

func (fls *FundListSpider) Run() {
	rm := NewResourceManagerChan(fls.threadsNum)
	for {
		url, ok := fls.scheduler.Pop()
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
	fls.parser.Parse(resp)

}

// 获取所有前端基金代号（后端基金跳过）
func (fls *FundListSpider) GetFundCodes() []string {
	return fls.parser.GetFundCodes()
}
