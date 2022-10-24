package main

import (
	. "COVID-19-Crawler/service"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	InitDB()
	for {
		Run()
		time.Sleep(4 * time.Hour)
	}
}

// init logger
func init() {
	//mkdir
	if _, err := os.Stat("tmp"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("tmp", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	//make file
	filename := "log.txt"
	logfile, err := os.Create("tmp/" + filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	out := io.MultiWriter(logfile, os.Stdout)
	log.SetOutput(out)
	log.SetPrefix("[Crawler] ")
}
