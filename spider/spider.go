package spider

import (
	"Fachoi_fund_test2/crawler"
	"Fachoi_fund_test2/scheduler"
)

type Spider struct {
	threadsNum int                       // 并发数
	urlsNum    int                       // 需要爬取的链接数
	scheduler  *scheduler.QueueScheduler // 队列调度器，用于管理url
	crawler    *crawler.Crawler          // 爬取器
	crawlCount chan int                  // 爬取进度统计
	//parseCount         chan int                  // 解析进度统计
	//saveCount          chan int                  // 存储进度统计
	//responseChan       chan interface{}          // http响应管道，用以将爬取器爬到的内容交给解析器解析
	//parsedResponseChan chan interface{}          // 经过解析的响应管道，用以将解析器解析后的内容交给存储器存储，每个Spider子类有不同格式的解析
	parser  interface{} // 解析器，由继承类确定类型
	saver   interface{} // 存储器，由继承类确定类型
	updater interface{}
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
	//s.parseCount = make(chan int, s.urlsNum)
	//s.saveCount = make(chan int, s.urlsNum)
	s.scheduler.Push(url)
}

func (s *Spider) AddUrls(urls []string) {
	s.urlsNum += len(urls)
	s.crawlCount = make(chan int, s.urlsNum)
	//s.parseCount = make(chan int, s.urlsNum)
	//s.saveCount = make(chan int, s.urlsNum)
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
