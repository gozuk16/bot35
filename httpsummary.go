package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func getPage(url string) (string, error) {
	var text string
	fmt.Println("URL(getPage): " + url)
	doc, err := goquery.NewDocument(url)
	doc.Find("head").Each(func(_ int, s *goquery.Selection) {
		title := s.Find("title").Text()
		fmt.Println("title(getPage): " + title)
		text += title
	})
	return text, err
}

func httpSummary(queryUrl string) (msg string) {
	res, _ := getPage(queryUrl)
	fmt.Println("URL: " + queryUrl)
	u, err := url.Parse(queryUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("host: " + u.Host)
	fmt.Println("Response: " + res)

	//msg = "> " + u.Host + "\n> <" + queryUrl + "|*" + res + "*>"
	msg = "> <" + queryUrl + "|*" + res + "*>"

	return
}
