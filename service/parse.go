package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

type Parser struct {
	sync.WaitGroup
	bodyData  []byte
	timestamp time.Time
}

func NewParser(bodyData []byte) *Parser {
	return &Parser{
		bodyData:  bodyData,
		timestamp: time.Now().UTC(),
	}
}

func saveAsJson(name string, data string) {
	//make dir
	if _, err := os.Stat("tmp"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("tmp", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	//make file
	filename := name + ".json"
	dstFile, err := os.Create("tmp/" + filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(dstFile)

	_, err = dstFile.WriteString(data + "\n")
	if err != nil {
		log.Println(err)
	}
}

func (p *Parser) Run() {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(p.bodyData))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Document length =", len(doc.Text()))

	//solve overall
	p.Add(1)
	go func() {
		overallInformation := regexp.MustCompile("({\"id\".*})}").FindString(doc.Find("script[id=getStatisticsService]").Text())
		if overallInformation != "" {
			p.overallParser(overallInformation)
		}
		p.Done()
	}()

	//solve area
	p.Add(1)
	go func() {
		areaInformation := regexp.MustCompile("\\[(.*)]").FindString(doc.Find("script[id=getAreaStat]").Text())
		if areaInformation != "" {
			p.provinceParser(areaInformation)
		}
		p.Done()

	}()

	//solve abroad
	p.Add(1)
	go func() {
		abroadInformation := regexp.MustCompile("\\[(.*)]").FindString(doc.Find("script[id=getListByCountryTypeService2true]").Text())
		if abroadInformation != "" {
			p.abroadParser(abroadInformation)
		}
		p.Done()
	}()

	//solve timelines
	p.Add(1)
	go func() {
		timelines := regexp.MustCompile(`\[(.*?)]`).FindString(doc.Find("script[id=getTimelineService1]").Text())
		p.timelinesParser(timelines)
		p.Done()
	}()

	p.Wait()
	log.Println("Successfully crawled!")
}

func (p *Parser) overallParser(data string) {
	log.Println("OverallParser running, len =", len(data))
	data = data[:len(data)-1]
	overallInfo := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &overallInfo)
	if err != nil {
		log.Fatal("[overallParser]", err)
	}
	delete(overallInfo, "id")
	delete(overallInfo, "createTime")
	delete(overallInfo, "modifyTime")
	delete(overallInfo, "imgUrl")
	delete(overallInfo, "deleted")
	overallInfo["updateTime"] = p.timestamp
	//?????????
	_, err = db.Collection("Overall").InsertOne(context.TODO(), overallInfo)
	if err != nil {
		log.Fatal(err)
	}
	saveAsJson("overallInfo", data)
}

func (p *Parser) provinceParser(data string) {
	log.Println("AreaParser running, len =", len(data))
	var areas []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &areas); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(areas); i++ {
		areas[i]["updateTime"] = p.timestamp
		_, err := db.Collection("Province").InsertOne(context.TODO(), areas[i])
		if err != nil {
			return
		}
	}

	saveAsJson("area", data)
}

func (p *Parser) abroadParser(data string) {
	log.Println("AbroadParser running, len =", len(data))
	var countries []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &countries); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(countries); i++ {
		delete(countries[i], "id")
		delete(countries[i], "tags")
		delete(countries[i], "sort")
		delete(countries[i], "modifyTime")
		delete(countries[i], "createTime")
		delete(countries[i], "countryType")
		delete(countries[i], "provinceId")
		delete(countries[i], "cityName")
		delete(countries[i], "provinceShortName")

		countries[i]["updateTime"] = p.timestamp
		//TODO INSERT
		_, err := db.Collection("Abroad").InsertOne(context.TODO(), countries[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	saveAsJson("abroad", data)
}

func (p *Parser) timelinesParser(data string) {
	log.Println("TimelineParser running, len =", len(data))
	var timelines []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &timelines); err != nil {
		log.Fatal(err)
	}
	//TODO INSERT
	for i := 0; i < len(timelines); i++ {
		timelines[i]["updateTime"] = p.timestamp
		_, err := db.Collection("Timeline").InsertOne(context.TODO(), timelines[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	saveAsJson("timelines", data)
}
