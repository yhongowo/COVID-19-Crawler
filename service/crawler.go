package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Crawler struct {
	sync.WaitGroup
}

func (c *Crawler) Run() {
	timestap := time.Now().Unix()
	doc := c.getDoc()

	//overall info
	go func() {
		c.Add(1)
		overallInformation := regexp.MustCompile("(\\{\"id\".*\\})\\}").FindString(doc.Find("script[id=getStatisticsService]").Text())
		if overallInformation != "" {
			c.overallParser(overallInformation)
		}
		c.Done()
	}()

	//area
	go func() {
		c.Add(1)
		areaInformation := regexp.MustCompile("\\[(.*)\\]").FindString(doc.Find("script[id=getStatisticsService]").Text())
		if areaInformation != "" {
			c.areaParser(areaInformation)
		}
		c.Done()
	}()

	//abroad
	go func() {
		c.Add(1)
		abroadInformation := regexp.MustCompile("\\[(.*)\\]").FindString(doc.Find("script[id=getAreaStat]").Text())
		if abroadInformation != "" {
			c.abroadParser(abroadInformation)
		}
		c.Done()
	}()

	//news chinese
	go func() {
		c.Add(1)
		newsChinese := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getTimelineService1]").Text())
		if newsChinese != "" {
			c.newsParser(newsChinese)
		}
		c.Done()
	}()

	//news english
	go func() {
		c.Add(1)
		newsEnglish := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getTimelineService2]").Text())
		if newsEnglish != "" {
			c.newsParser(newsEnglish)
		}
		c.Done()
	}()

	//rumors
	go func() {
		c.Add(1)
		rumors := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getIndexRumorList]").Text())
		if rumors != "" {
			c.rumorParser(rumors)
		}
		c.Done()
	}()

	c.Wait()
	fmt.Println("[INFO]:Successfully crawled!")
}

func (c *Crawler) getDoc() *goquery.Document {
	client := http.Client{}
	req, err := http.NewRequest("Get", "https://ncov.dxy.cn/ncovh5/view/pneumonia", nil)
	if err != nil {
		fmt.Println(err)
	}
	//add header
	req.Header.Add("User-Agent", Agent())
	//send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	//defer close
	defer resp.Body.Close()
	//read resp
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return doc
}

func (c *Crawler) overallParser(info string) {
	
}

func (c *Crawler) areaParser(info string) {

}

func (c *Crawler) abroadParser(info string) {

}

func (c *Crawler) newsParser(info string) {

}

func (c *Crawler) rumorParser(info string) {

}
