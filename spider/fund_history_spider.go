package spider

import (
	"Fachoi_fund_test/parser"
	"Fachoi_fund_test/resource_manager"
	"Fachoi_fund_test/saver"
	"Fachoi_fund_test/update_checker"
	"fmt"
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
	"time"
)

type FundHistorySpider struct {
	*Spider
	parser        *parser.FundHistoryParser
	saver         *saver.FundHistorySaver
	updateChecker *update_checker.FundHistoryUpdateChecker
}

func NewFundHistorySpider(db *sqlx.DB, threadsNum int) *FundHistorySpider {
	fhs := new(FundHistorySpider)
	fhs.Spider = NewSpider(threadsNum)
	fhs.parser = parser.NewFundHistoryParser()
	fhs.saver = saver.NewFundHistorySaver(db)
	fhs.updateChecker = update_checker.NewFundHistoryUpdateChecker(db)
	return fhs
}

func (fhs *FundHistorySpider) Run() {
	crawlerThreadsManager := resource_manager.NewResourceManager(fhs.threadsNum)
	for {
		url, ok := fhs.scheduler.Pop()
		if ok == false && crawlerThreadsManager.Has() == 0 {
			fmt.Println("基金历史数据存储完毕！")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		crawlerThreadsManager.GetOne()
		go func(url string) {
			if url == "" {
				return
			}
			defer crawlerThreadsManager.FreeOne()
			fhs.process(url)
		}(url)
	}
}

func (fhs *FundHistorySpider) process(url string) {
	resp := fhs.crawler.Crawl(url)
	if resp == nil {
		fmt.Println("resp is nil")
		if !fhs.crawler.UrlVisited(url) {
			fhs.scheduler.Push(url)
		}
	} else {
		// 识别出url中的基金代码
		reg := regexp.MustCompile(`fundCode=\d+`)
		code := strings.Split(reg.FindAllString(url, -1)[0], "=")[1]
		parsedResp := fhs.parser.Parse(resp)
		fhs.saver.Save(parsedResp, code)
		fhs.crawlCount <- 1
		fmt.Printf("爬取进度：%d / %d \n", len(fhs.crawlCount), fhs.urlsNum)
	}
}

func (fhs *FundHistorySpider) Update() {
	crawlerThreadsManager := resource_manager.NewResourceManager(fhs.threadsNum)
	for {
		url, ok := fhs.scheduler.Pop()
		if ok == false && crawlerThreadsManager.Has() == 0 {
			fmt.Println("基金历史数据爬取存储完毕！")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		crawlerThreadsManager.GetOne()
		go func(url string) {
			if url == "" {
				return
			}
			defer crawlerThreadsManager.FreeOne()
			fhs.updateProcess(url)
		}(url)
	}
}

func (fhs *FundHistorySpider) updateProcess(url string) {
	resp := fhs.crawler.Crawl(url)
	if resp == nil {
		fmt.Println("resp is nil")
		if !fhs.crawler.UrlVisited(url) {
			fhs.scheduler.Push(url)
		}
	} else {
		// 识别出url中的基金代码
		reg := regexp.MustCompile(`fundCode=\d+`)
		code := strings.Split(reg.FindAllString(url, -1)[0], "=")[1]
		parsedResp := fhs.parser.Parse(resp)
		candidateResp := fhs.updateChecker.Check(parsedResp, code)
		fhs.saver.Save(candidateResp, code)
		fhs.crawlCount <- 1
		fmt.Printf("爬取进度：%d / %d \n", len(fhs.crawlCount), fhs.urlsNum)
	}
}
