package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ConfluenceQuestions []struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Author struct {
		Name     string `json:"name"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		UserKey  string `json:"userKey"`
	} `json:"author"`
	FriendlyDateAsked string `json:"friendlyDateAsked"`
	DateAsked         int64  `json:"dateAsked"`
	AnswersCount      int    `json:"answersCount"`
	Topics            []struct {
		ID         int    `json:"id"`
		IDAsString string `json:"idAsString"`
		Name       string `json:"name"`
		URL        string `json:"url"`
		Featured   bool   `json:"featured"`
		IsWatching bool   `json:"isWatching"`
	} `json:"topics"`
	Votes struct {
		Up        int  `json:"up"`
		Down      int  `json:"down"`
		Total     int  `json:"total"`
		UpVoted   bool `json:"upVoted"`
		DownVoted bool `json:"downVoted"`
	} `json:"votes"`
	Space struct {
		SpaceKey  string `json:"spaceKey"`
		SpaceName string `json:"spaceName"`
	} `json:"space"`
	IsTrashed bool `json:"isTrashed"`
}

var confluenceQuestions ConfluenceQuestions

func encodeJson4Confluence(url string) (ConfluenceQuestions, error) {
	m := map[string]string{"url": url}
	resp, err := getResponse(m)
	if err != nil {
		log.Fatalf("Error!: %v", err)
		return confluenceQuestions, err
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error!: %v", err)
		return confluenceQuestions, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(byteArray, &confluenceQuestions); err != nil {
		log.Fatalf("Error!: %v", err)
	}

	fmt.Printf("%v\n", byteArray)
	return confluenceQuestions, err
}

func confluence(url string) (msg string) {
	res, _ := encodeJson4Confluence(url)
	log.Printf("unanswerd=%d\n", len(res))
	for _, v := range res {
		fmt.Printf("v.Title: %v\n", v.Title)
		msg += v.Title + "/"
		fmt.Printf("v.Author.Name: %v\n", v.Author.Name)
		msg += v.Author.Name + " "
		fmt.Printf("v.URL: %v\n", v.URL)
		msg += "[" + v.URL + "]\n"
	}

	/*
		fmt.Println("Project: " + res.Fields.Project.Name)
		fmt.Println("Title: " + res.Fields.Summary)
		fmt.Println("Status: " + res.Fields.Status.Name)
		fmt.Println("Issuetype: " + res.Fields.Issuetype.Name)
	*/

	//ticketUrl := config.Confluence.Url + res.Key
	//msg = res.Key + " (" + res.Fields.Status.Name + ") " + res.Fields.Project.Name + " : " + res.Fields.Summary + " [" + ticketUrl + "]"
	msg = ""

	return
}
