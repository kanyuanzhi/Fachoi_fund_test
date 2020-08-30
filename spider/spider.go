package spider

import (
	"Fachoi_fund_test2/crawler"
	"Fachoi_fund_test2/scheduler"
	"fmt"
	"time"
)

type Spider struct {
	threadsNum uint8
	scheduler  *scheduler.QueueScheduler
	crawler    *crawler.Crawler
	parser     interface{}
	saver      interface{}
}

func NewSpider(threadsNum uint8) *Spider {
	return &Spider{
		threadsNum: threadsNum,
		scheduler:  scheduler.NewQueueScheduler(),
		crawler:    crawler.NewCrawler(),
	}
}

func (s *Spider) AddUrl(url string) {
	s.scheduler.Push(url)
}

func (s *Spider) AddUrls(urls []string) {
	for _, url := range urls {
		s.scheduler.Push(url)
	}
}

func (s *Spider) Run() {
	rm := NewResourceManagerChan(s.threadsNum)
	for {
		url, ok := s.scheduler.Pop()
		if ok == false && rm.Has() == 0 {
			fmt.Println("url爬取完毕")
			break
		} else if ok == false {
			time.Sleep(time.Second)
			continue
		}
		rm.GetOne()
		go func(url string) {
			defer rm.FreeOne()
			s.Process(url)
		}(url)
	}
}

func (s *Spider) Process(url string) {
}
