package resource_manager

// 爬取器资源管理器，用于控制并发爬取数量
type CrawlerResourceManager struct {
	*ResourceManager
}

func NewCrawlerResourceManager(num int) *CrawlerResourceManager {
	crm := new(CrawlerResourceManager)
	crm.ResourceManager = NewResourceManager(num)
	return crm
}
