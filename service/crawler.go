package service

import (
	"github.com/gocolly/colly/v2"
	"log"
)

const WebsiteURI = "https://ncov.dxy.cn/ncovh5/view/pneumonia"

// Run Crawler
func Run() {
	//Initialize Collector
	c := colly.NewCollector(colly.AllowedDomains())

	//Register callback function
	c.OnRequest(func(request *colly.Request) {
		log.Println("Start Crawler...")
	})
	c.OnResponse(func(response *colly.Response) {
		NewParser(response.Body).Run()
	})
	c.OnError(func(response *colly.Response, err error) {
		log.Println("[ERROR]", err)
	})
	c.OnScraped(func(response *colly.Response) {
		log.Println("Stop Crawler...")
	})
	//Start Collector
	if err := c.Visit(WebsiteURI); err != nil {
		log.Println(err)
	}
}
