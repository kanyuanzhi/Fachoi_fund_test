package spider

import (
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/saver"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type FundListSpider struct {
	*Spider
	parser *parser.FundListParser
	saver  *saver.FundListSaver
}

func NewFundListSpider(db *sqlx.DB) *FundListSpider {
	fls := new(FundListSpider)
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	fls.saver = saver.NewFundListSaver(db)
	return fls
}

func (fls *FundListSpider) Run() {
	for {
		url, ok := fls.scheduler.Pop()
		if ok == false {
			break
		}
		fls.process(url)
	}
	codes := fls.GetFrontFundCodes()
	fmt.Printf("基金列表爬取完毕！共爬取%d个基金，其中前端基金%d个。\n", fls.urlsNum, len(codes))
	//util.CreateAllFundHistoryTables(fls.db, codes)
}

func (fls *FundListSpider) process(url string) {
	resp := fls.crawler.Crawl(url)
	if resp == nil {
		if !fls.crawler.UrlVisited(url) {
			fls.scheduler.Push(url)
		}
		return
	}
	parsedResp := fls.parser.Parse(resp)
	fls.saver.Save(parsedResp)
}

// 获取所有前端基金代号（后端基金跳过）
func (fls *FundListSpider) GetFrontFundCodes() []string {
	return fls.parser.GetFrontFundCodes()
}
