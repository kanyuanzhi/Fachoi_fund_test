package spider

import (
	"Fachoi_fund_test2/crawler"
	"Fachoi_fund_test2/scheduler"
)

type Spider struct {
	threadsNum int
	urlsNum    int
	scheduler  *scheduler.QueueScheduler
	crawler    *crawler.Crawler
	parser     interface{}
	saver      interface{}
	crawlCount chan int
	saveCount  chan int
}

func NewSpider(threadsNum int) *Spider {
	return &Spider{
		threadsNum: threadsNum,
		urlsNum:    0,
		scheduler:  scheduler.NewQueueScheduler(),
		crawler:    crawler.NewCrawler(),
	}
}

func (s *Spider) AddUrl(url string) {
	s.urlsNum += 1
	s.crawlCount = make(chan int, s.urlsNum)
	s.saveCount = make(chan int, s.urlsNum)
	s.scheduler.Push(url)
}

func (s *Spider) AddUrls(urls []string) {
	s.urlsNum += len(urls)
	s.crawlCount = make(chan int, s.urlsNum)
	s.saveCount = make(chan int, s.urlsNum)
	for _, url := range urls {
		s.scheduler.Push(url)
	}
}

//func Run(s *FundInfoSpider) {
//	rm := resource_manager.NewResourceManager(s.threadsNum)
//	for {
//		url, ok := s.scheduler.Pop()
//		if ok == false && rm.Has() == 0 {
//			fmt.Println("url爬取完毕")
//			break
//		} else if ok == false {
//			time.Sleep(time.Second)
//			continue
//		}
//		rm.GetOne()
//		go func(url string) {
//			defer rm.FreeOne()
//			s.process(url)
//		}(url)
//	}
//}

//func (s *Spider) Run() {
//	rm := NewResourceManagerChan(s.threadsNum)
//	for {
//		url, ok := s.scheduler.Pop()
//		if ok == false && rm.Has() == 0 {
//			fmt.Println("url爬取完毕")
//			break
//		} else if ok == false {
//			time.Sleep(time.Second)
//			continue
//		}
//		rm.GetOne()
//		go func(url string) {
//			defer rm.FreeOne()
//			s.Process(url)
//		}(url)
//	}
//}
//
//func (s *Spider) Process(url string) {
//}
