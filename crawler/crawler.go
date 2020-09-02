package crawler

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
)

type Crawler struct {
	locker  *sync.Mutex
	crawled map[[md5.Size]byte]bool
}

func NewCrawler() *Crawler {
	locker := new(sync.Mutex)
	crawled := make(map[[md5.Size]byte]bool)
	return &Crawler{
		locker:  locker,
		crawled: crawled,
	}
}

func (c *Crawler) Crawl(url string) *http.Response {
	key := md5.Sum([]byte(url))

	client := &http.Client{}
	req, _ := http.NewRequest("Get", url, nil)
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"
	referer := "http://fund.eastmoney.com/allfund.html"
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", referer)
	resp, err := client.Do(req)

	c.locker.Lock()
	if has, ok := c.crawled[key]; has && ok {
		c.locker.Unlock()
		return nil
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("爬取链接失败：", url)
		c.crawled[key] = false
		c.locker.Unlock()
		return nil
	}

	c.crawled[key] = true
	c.locker.Unlock()
	return resp
}

func (c *Crawler) UrlVisited(url string) bool {
	key := md5.Sum([]byte(url))
	var visited bool
	c.locker.Lock()
	if ok, has := c.crawled[key]; ok && has {
		visited = true
	} else {
		visited = false
	}
	c.locker.Unlock()
	return visited
}
