package spider

import (
	"Fachoi_fund_test2/parser"
)

type FundListSpider struct {
	*Spider
	parser *parser.FundListParser
}

func NewFundListSpider() *FundListSpider {
	var fls *FundListSpider
	fls.Spider = NewSpider(1)
	fls.parser = parser.NewFundListParser()
	return fls
}
