package service

import (
	"COVID-19-Crawler/util"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

type Parser struct {
	sync.WaitGroup
	bodyData  []byte
	timestamp int64
}

func NewParser(bodyData []byte) *Parser {
	return &Parser{
		bodyData:  bodyData,
		timestamp: time.Now().Unix(),
	}
}

func SaveAsJson(name string, data string) {
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
}

func (p *Parser) overallParser(data string) {
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
	//持久化
	filter := bson.D{{"updateTime", bson.D{{"$gte", util.TodayBeginTime()}}}}
	opts := options.Replace().SetUpsert(true)
	_, err = db.Collection("Overall").ReplaceOne(context.TODO(), filter, overallInfo, opts)
	if err != nil {
		log.Fatal(err)
	}
	SaveAsJson("overallInfo", data)
	log.Println("Save overall info")
}

func (p *Parser) provinceParser(data string) {
	var areas []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &areas); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(areas); i++ {
		areas[i]["updateTime"] = p.timestamp
		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"updateTime", bson.D{{"$gte", util.TodayBeginTime()}}}},
					bson.D{{"provinceName", areas[i]["provinceName"]}}},
			},
		}
		opts := options.Replace().SetUpsert(true)
		_, err := db.Collection("Area").ReplaceOne(context.TODO(), filter, areas[i], opts)
		if err != nil {
			return
		}
	}
	SaveAsJson("area", data)
	log.Printf("Save province info,rows:%d", len(areas))

}

func (p *Parser) abroadParser(data string) {
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

		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"updateTime", bson.D{{"$gte", util.TodayBeginTime()}}}},
					bson.D{{"provinceName", countries[i]["provinceName"]}}},
			},
		}
		opts := options.Replace().SetUpsert(true)
		_, err := db.Collection("Abroad").ReplaceOne(context.TODO(), filter, countries[i], opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	SaveAsJson("abroad", data)
	log.Printf("Save abroad info,rows:%d", len(countries))
}

func (p *Parser) timelinesParser(data string) {
	var timelines []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &timelines); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(timelines); i++ {
		timelines[i]["updateTime"] = p.timestamp

		filter := bson.D{{"articleId", timelines[i]["articleId"]}}
		opts := options.Replace().SetUpsert(true)
		_, err := db.Collection("Timeline").ReplaceOne(context.TODO(), filter, timelines[i], opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	SaveAsJson("timelines", data)
	log.Printf("Save timeline info,rows:%d", len(timelines))
}
