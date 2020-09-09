package spider

import (
	"Fachoi_fund_test/parser"
	"Fachoi_fund_test/resource_manager"
	"Fachoi_fund_test/saver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type FundInfoSpider struct {
	*Spider
	parser *parser.FundInfoParser
	saver  *saver.FundInfoSaver
}

func NewFundInfoSpider(db *sqlx.DB, threadsNum int) *FundInfoSpider {
	fis := new(FundInfoSpider)
	fis.Spider = NewSpider(threadsNum)
	fis.parser = parser.NewFundInfoParser()
	fis.saver = saver.NewFundInfoSaver(db)
	return fis
}

func (fis *FundInfoSpider) Run() {
	crm := resource_manager.NewResourceManager(fis.threadsNum)
	//util.TruncateTable("fund_info_table", fis.saver.DB)
	for {
		url, ok := fis.scheduler.Pop()
		if ok == false && crm.Has() == 0 {
			fmt.Println("基金信息爬取存储完毕！")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		crm.GetOne()
		// 并发爬取并解析页面
		go func(url string) {
			if url == "" {
				return
			}
			defer crm.FreeOne()
			fis.process(url)
		}(url)
	}
}

func (fis *FundInfoSpider) process(url string) {
	resp := fis.crawler.Crawl(url)
	if resp == nil {
		if !fis.crawler.UrlVisited(url) {
			fis.scheduler.Push(url)
		}
	} else {
		parsedResp := fis.parser.Parse(resp)
		fis.saver.Save(parsedResp)
		fis.crawlCount <- 1
		fmt.Printf("爬取进度：%d / %d \n", len(fis.crawlCount), fis.urlsNum)
	}
}

func (fis *FundInfoSpider) Update() {
	fis.Run()
}
