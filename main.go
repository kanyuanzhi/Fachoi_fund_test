package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	fmt.Println("main")
	url := "http://fund.eastmoney.com/js/fundcode_search.js"
	key := md5.Sum([]byte(url))
	fmt.Println(md5.Size)
	fmt.Println(key)
	fmt.Printf("%x\n", key)

	//c := crawler.NewCrawler()
	//resp := c.Crawl(url)
	//
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//bodyStr := string(bodyBytes)

	fmt.Println([]byte(url))

}
