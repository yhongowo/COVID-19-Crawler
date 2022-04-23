package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Run() {
	timestap := time.Now().Unix()
	doc := c.getDoc()

	//overall info
	overallInformation := regexp.MustCompile("(\\{\"id\".*\\})\\}").FindString(doc.Find("script[id=getStatisticsService]").Text())
	if overallInformation != "" {
		c.overallParser(overallInformation)
	}

	//area
	areaInformation := regexp.MustCompile("\\[(.*)\\]").FindString(doc.Find("script[id=getStatisticsService]").Text())
	if areaInformation != "" {
		c.areaParser(areaInformation)
	}

	//abroad
	abroadInformation := regexp.MustCompile("\\[(.*)\\]").FindString(doc.Find("script[id=getAreaStat]").Text())
	if abroadInformation != "" {
		c.abroadParser(abroadInformation)
	}

	//news chinese
	newsChinese := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getTimelineService1]").Text())
	if newsChinese != "" {
		c.newsParser(newsChinese)
	}

	//news english
	newsEnglish := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getTimelineService2]").Text())
	if newsEnglish != "" {
		c.newsParser(newsEnglish)
	}

	//rumors
	rumors := regexp.MustCompile("\\[(.*?)\\]").FindString(doc.Find("script[id=getIndexRumorList]").Text())
	if rumors != "" {
		c.rumorParser(rumors)
	}
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
