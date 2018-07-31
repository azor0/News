package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
)
// type Location struct {
// 	Loc string `xml:"loc"`
// }
  
// // Why doesnt this work
// func (l Location) String() string {
// 	fmt.Sprintf(l.Loc)
// 	httpsFix := l.Loc[:4] + "s" + l.Loc[4:]
// 	return httpsFix
// }



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
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
}
// func makeRequest(url string)

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex
	var n News

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.washingtonpost.com/news-sitemap-index.xml?", nil)
	cookie1 := http.Cookie{Name: "wp_gdpr", Value: "1|1" }
    req.AddCookie(&cookie1)
	res, _ := client.Do(req)

	bytes, _ := ioutil.ReadAll(res.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap) // key is string and values are NewsMap values


	for _, Location := range s.Locations {
		testing := fmt.Sprintf("%s s %s", Location[:4], Location[4:])
		url := Location[:4] + "s" + Location[4:]
		req, _ := http.NewRequest("GET", url, nil)
		cookie1 := http.Cookie{Name: "wp_gdpr", Value: "1|1" }
		req.AddCookie(&cookie1)
		fmt.Println("testing:", testing)
		test, _ := client.Do(req)

		bytes, _ := ioutil.ReadAll(test.Body)
		string_body := string(bytes)
		fmt.Println(string_body)
		xml.Unmarshal(bytes, &n)
		for idx, _ := range n.Keywords {
			fmt.Println("Heere2")
			fmt.Println(n.Titles[idx], idx)
			news_map[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
	}

	p := NewsAggPage{Title: "NEWS", News: news_map}
    t, _ := template.ParseFiles("template.html")
    t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":8000", nil) 
}