package main

import (
	"COVID-19-Crawler/service"
	"time"
)

func main() {
	//Schedule task
	for true {
		service.InitDB()
		service.Run()
		service.CloseDB()
		//Set duration here
		time.Sleep(1 * time.Hour)
	}
}
