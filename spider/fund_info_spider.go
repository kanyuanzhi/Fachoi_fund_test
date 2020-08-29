package spider

type FundInfoSpider struct {
	Spider
}

func NewFundInfoSpider(threadsNum uint) *FundInfoSpider {
	return &FundInfoSpider{}
}
