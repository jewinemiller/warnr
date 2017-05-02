package main

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"os"
	"time"
	"log"
	"github.com/PuerkitoBio/goquery"
	tm "github.com/buger/goterm"
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

	feed, _ := fp.ParseURL("https://alerts.weather.gov/cap/us.php?x=0")
	parseFeed(feed, topic, eventMap)
	
	for range time.Tick(time.Minute * 1){
		feed, _ := fp.ParseURL("https://alerts.weather.gov/cap/us.php?x=0")
		parseFeed(feed, topic, eventMap)	
	}
}

func parseFeed(_feed *gofeed.Feed, topic string, eventMap map[string]*gofeed.Item){
	tm.Clear()
	tm.MoveCursor(1,1)
	printedCount := 0
	for i:=0; i<len(_feed.Items); i++{
		event := _feed.Items[i].Extensions["cap"]["event"]
		if topic == "" || strings.Contains(strings.ToLower(event[0].Value), strings.ToLower(topic)){
			//if _, ok := eventMap[_feed.Items[i].GUID]; !ok {
			parseLink(_feed.Items[i].GUID)
			printedCount++
				//eventMap[_feed.Items[i].GUID] = _feed.Items[i]
				//expiration := _feed.Items[i].Extensions["cap"]["expires"]
				//effective := _feed.Items[i].Extensions["cap"]["effective"]
				//fmt.Println("Event: ", event[0].Value)
				//fmt.Println("Effective: ", effective[0].Value)
				//fmt.Println("Expires: ", expiration[0].Value)
				//fmt.Println("Description: ", _feed.Items[i].Description)
				//fmt.Println("--------------------------\n")  
			//}
		}
	}
	if printedCount == 0{
		tm.Println("No Active Alerts")
	}

	tm.Flush()
}

func parseLink(link string){
	doc, err := goquery.NewDocument(link)
	if err != nil {
	  log.Fatal(err)
	}

	doc.Find("headline").Each(func(i int, s *goquery.Selection){
		headline := s.Text()
		tm.Println(headline)	
	})

	doc.Find("effective").Each(func(i int, s *goquery.Selection){
		effective := s.Text()
		tm.Println("Effective: ", effective)	
	})

	doc.Find("expires").Each(func(i int, s *goquery.Selection){
		expires := s.Text()
		tm.Println("Expires: ", expires)	
	})

	// Find the review items
	doc.Find("description").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
	    description := s.Text()
	    tm.Println("DESCRIPTION:", description)
	 })

	tm.Println("-------------------------------------------------------------\n")
}

