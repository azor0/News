package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
	"sync"
)

var wg sync.WaitGroup
var client http.Client

type NewsMap struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
    Title string
    News map[string]NewsMap
}

type Sitemapindex struct {
	Locations []string `xml:"sitemap>loc"`
}


type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Go to /news</h1>")
}

func newsRoutine(c chan News, Location string){
	defer wg.Done()

	var n News
	url := Location[:4] + "s" + Location[4:]
	req, _ := http.NewRequest("GET", url, nil)
	cookie1 := http.Cookie{Name: "wp_gdpr", Value: "1|1" }
	req.AddCookie(&cookie1)
	response, _ := client.Do(req)

	bytes, _ := ioutil.ReadAll(response.Body)
	xml.Unmarshal(bytes, &n)
	response.Body.Close()

	c <- n
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex

	req, _ := http.NewRequest("GET", "https://www.washingtonpost.com/news-sitemap-index.xml?", nil)
	cookie1 := http.Cookie{Name: "wp_gdpr", Value: "1|1" }
    req.AddCookie(&cookie1)
	res, _ := client.Do(req)

	bytes, _ := ioutil.ReadAll(res.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)
	res.Body.Close()
	queue := make(chan News, 30)

	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}
	wg.Wait()
	close(queue)
	for elem := range queue { 
		for idx, _ := range elem.Keywords {
			news_map[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
		}
	}

	p := NewsAggPage{Title: "NEWS WEBSITE", News: news_map}
    t, _ := template.ParseFiles("template.html")
    t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/news/", newsHandler)
	http.ListenAndServe(":8000", nil) 
}