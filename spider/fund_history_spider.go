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
	parser             *parser.FundHistoryParser
	saver              *saver.FundHistorySaver
	responseChan       chan *ResponseAndCode
	parsedResponseChan chan *ParsedResponseAndCode
}

type ResponseAndCode struct {
	Resp *http.Response
	Code string
}

func NewResponseAndCode(Resp *http.Response, Code string) *ResponseAndCode {
	return &ResponseAndCode{
		Resp: Resp,
		Code: Code,
	}
}

type ParsedResponseAndCode struct {
	ParsedResp []db_model.FundHistoryModel
	Code       string
}

func NewParsedResponseAndCode(ParsedResp []db_model.FundHistoryModel, Code string) *ParsedResponseAndCode {
	return &ParsedResponseAndCode{
		ParsedResp: ParsedResp,
		Code:       Code,
	}
}

func NewFundHistorySpider(db *sql.DB, threadsNum int) *FundHistorySpider {
	fhs := new(FundHistorySpider)
	fhs.Spider = NewSpider(threadsNum)
	fhs.parser = parser.NewFundHistoryParser()
	fhs.saver = saver.NewFundHistorySaver(db)
	fhs.responseChan = make(chan *ResponseAndCode, 10000)
	fhs.parsedResponseChan = make(chan *ParsedResponseAndCode, 10000)
	return fhs
}

func (fhs *FundHistorySpider) Run() {

	crawlerThreadsManager := resource_manager.NewResourceManager(fhs.threadsNum)
	saverThreadsManager := resource_manager.NewResourceManager(fhs.threadsNum)
	saverProcessManager := resource_manager.NewResourceManager(fhs.urlsNum)
	saverProcessManager.FillToTheFull()
	crawlFinished := false

	for {
		url, ok := fhs.scheduler.Pop()
		if crawlFinished == false && ok == false && crawlerThreadsManager.Has() == 0 {
			fmt.Println("基金历史数据爬取完毕！等待存储完毕......")
			crawlFinished = true
		} else if crawlFinished == true && ok == false && saverProcessManager.Has() == 0 {
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
			fhs.crawlProcess(url)
		}(url)

		go func() {
			fhs.parseProcess()
		}()

		saverThreadsManager.GetOne()
		go func() {
			defer func() {
				saverThreadsManager.FreeOne()
				saverProcessManager.FreeOne()
			}()
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
		// 识别出url中的基金代码
		reg := regexp.MustCompile(`fundCode=\d+`)
		code := strings.Split(reg.FindAllString(url, -1)[0], "=")[1]
		fhs.responseChan <- NewResponseAndCode(resp, code)
		fhs.crawlCount <- 1
		fmt.Printf("爬取进度：%d / %d \n", len(fhs.crawlCount), fhs.urlsNum)
	}
}

func (fhs *FundHistorySpider) parseProcess() {
	respAndCode := <-fhs.responseChan
	resp := respAndCode.Resp
	code := respAndCode.Code
	parsedResp := fhs.parser.Parse(resp)
	fhs.parsedResponseChan <- NewParsedResponseAndCode(parsedResp, code)
	fhs.parseCount <- 1
	fmt.Printf("解析进度：%d / %d \n", len(fhs.parseCount), fhs.urlsNum)
}

func (fhs *FundHistorySpider) saveProcess() {
	parsedRespAndCode := <-fhs.parsedResponseChan
	parsedResp := parsedRespAndCode.ParsedResp
	code := parsedRespAndCode.Code
	fhs.saver.Save(parsedResp, code)
	fhs.saveCount <- 1
	fmt.Printf("存储进度：%d / %d \n", len(fhs.saveCount), fhs.urlsNum)
}
