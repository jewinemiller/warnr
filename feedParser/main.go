package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
	"os"
	"time"
)

func main(){
	var topic string
	if(len(os.Args) < 2){
		topic = ""
	}else{
		topic = os.Args[1]
	}
	var eventMap map[string]*gofeed.Item = make(map[string]*gofeed.Item)
	fp := gofeed.NewParser()

	//Initial run
	feed, _ := fp.ParseURL("https://alerts.weather.gov/cap/us.php?x=0")
	parseFeed(feed, topic, eventMap)	

	for range time.Tick(time.Minute * 1){
		feed, _ := fp.ParseURL("https://alerts.weather.gov/cap/us.php?x=0")
		parseFeed(feed, topic, eventMap)	
	}
}

func parseFeed(_feed *gofeed.Feed, topic string, eventMap map[string]*gofeed.Item){
	for i:=0; i<len(_feed.Items); i++{
		event := _feed.Items[i].Extensions["cap"]["event"]
		if topic == "" || strings.Contains(strings.ToLower(event[0].Value), strings.ToLower(topic)){
			if _, ok := eventMap[_feed.Items[i].GUID]; !ok {
				eventMap[_feed.Items[i].GUID] = _feed.Items[i]
				expiration := _feed.Items[i].Extensions["cap"]["expires"]
				effective := _feed.Items[i].Extensions["cap"]["effective"]
				fmt.Println("Event: ", event[0].Value)
				fmt.Println("Effective: ", effective[0].Value)
				fmt.Println("Expires: ", expiration[0].Value)
				fmt.Println("Description: ", _feed.Items[i].Description)
				fmt.Println("--------------------------\n")  
			}
		}
	}
}
