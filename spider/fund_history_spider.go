package spider

import (
	"Fachoi_fund_test2/db_model"
	"Fachoi_fund_test2/parser"
	"Fachoi_fund_test2/resource_manager"
	"Fachoi_fund_test2/saver"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type FundHistorySpider struct {
	*Spider
	parser       *parser.FundHistoryParser
	saver        *saver.FundHistorySaver
	responseChan chan *http.Response
	//parsedResponseChan chan db_model.FundHistoryModelAndCode
	parsedResponseChan chan []db_model.FundHistoryModel
	codeChan           chan string
}

func NewFundHistorySpider(db *sql.DB, threadsNum int) *FundHistorySpider {
	fhs := new(FundHistorySpider)
	fhs.Spider = NewSpider(threadsNum)
	fhs.parser = parser.NewFundHistoryParser()
	fhs.saver = saver.NewFundHistorySaver(db)
	fhs.responseChan = make(chan *http.Response, 1000)
	fhs.codeChan = make(chan string, 1000)
	fhs.parsedResponseChan = make(chan []db_model.FundHistoryModel, 1000)
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
		} else if crawlFinished && srm.Has() == 0 {
			fmt.Println("基金历史数据存储完毕！")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		crm.GetOne()
		go func(url string) {
			defer crm.FreeOne()
			fhs.crawlProcess(url)
		}(url)

		go func() {
			fhs.parseProcess()
		}()

		go func() {
			defer srm.FreeOne()
			fhs.saveProcess()
		}()
	}
}

func (fhs *FundHistorySpider) crawlProcess(url string) {
	resp := fhs.crawler.Crawl(url)
	if resp == nil {
		if !fhs.crawler.UrlVisited(url) {
			fhs.scheduler.Push(url)
		}
	} else {
		fhs.responseChan <- resp
		// 识别出url中的基金代码
		reg := regexp.MustCompile(`fundCode=\d+`)
		fhs.codeChan <- strings.Split(reg.FindAllString(url, -1)[0], "=")[1]
		fhs.crawlCount <- 1
		fmt.Printf("爬取进度：%d / %d \n", len(fhs.crawlCount), fhs.urlsNum)
	}
}

func (fhs *FundHistorySpider) parseProcess() {
	//fhs.parsedResponseChan <- db_model.FundHistoryModelAndCode{Fhms: fhs.parser.Parse(<-fhs.responseChan), Code: <-fhs.codeChan}
	fhs.parsedResponseChan <- fhs.parser.Parse(<-fhs.responseChan)
	fhs.parseCount <- 1
	fmt.Printf("解析进度：%d / %d \n", len(fhs.parseCount), fhs.urlsNum)
}

func (fhs *FundHistorySpider) saveProcess() {
	fhs.saver.Save(<-fhs.parsedResponseChan, <-fhs.codeChan)
	fhs.saveCount <- 1
	fmt.Printf("存储进度：%d / %d \n", len(fhs.saveCount), fhs.urlsNum)
}
