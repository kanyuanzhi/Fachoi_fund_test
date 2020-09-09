package spider

import (
	"Fachoi_fund_test/parser"
	"Fachoi_fund_test/saver"
	"Fachoi_fund_test/update_checker"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type FundListSpider struct {
	*Spider
	parser        *parser.FundListParser
	saver         *saver.FundListSaver
	updateChecker *update_checker.FundListUpdateChecker
}

func NewFundListSpider(db *sqlx.DB) *FundListSpider {
	fls := new(FundListSpider)
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	fls.saver = saver.NewFundListSaver(db)
	fls.updateChecker = update_checker.NewFundListUpdateChecker(db)
	return fls
}

func (fls *FundListSpider) Run() {
	//util.TruncateTable("fund_list_table", fls.saver.DB)
	for {
		url, ok := fls.scheduler.Pop()
		if ok == false {
			break
		}
		fls.process(url)
	}
	codes := fls.GetFrontFundCodes()
	fmt.Printf("基金列表爬取完毕！其中前端基金%d个。\n", len(codes))
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

func (fls *FundListSpider) Update() {
	for {
		url, ok := fls.scheduler.Pop()
		if ok == false {
			break
		}
		func() {
			resp := fls.crawler.Crawl(url)
			if resp == nil {
				if !fls.crawler.UrlVisited(url) {
					fls.scheduler.Push(url)
				}
				return
			}
			parsedResp := fls.parser.Parse(resp)
			candidateResp := fls.updateChecker.Check(parsedResp)
			if len(candidateResp) == 0 {
				fmt.Println("无新基金需要更新!")
				return
			}
			fls.saver.Save(candidateResp)
			fmt.Printf("更新完毕！共更新%d个基金，其中前端基金%d个\n", len(candidateResp), len(fls.GetUpdatedFrontFundCodes()))
		}()
	}
}

func (fls *FundListSpider) GetUpdatedFrontFundCodes() []string {
	return fls.updateChecker.GetUpdatedFrontFundCodes()
}
