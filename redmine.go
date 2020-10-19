package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"log"
	"strconv"
)

type RedmineIssue struct {
	Issue struct {
		ID      int `json:"id"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Tracker struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tracker"`
		Status struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"priority"`
		Author struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"author"`
		AssignedTo struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"assigned_to"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"category"`
		Subject      string `json:"subject"`
		Description  string `json:"description"`
		DoneRatio    int    `json:"done_ratio"`
		CustomFields []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"custom_fields"`
		CreatedOn time.Time `json:"created_on"`
		UpdatedOn time.Time `json:"updated_on"`
		ClosedOn  time.Time `json:"closed_on"`
	} `json:"issue"`
}

var redmineIssue RedmineIssue

func encodeJson4Redmine(url string) (RedmineIssue, error) {
	fmt.Println("headerValue" + config.Redmine.APIToken)
	m := map[string]string{
		"url":         url + ".json",
		"headerKey":   "X-Redmine-API-Key",
		"headerValue": config.Redmine.APIToken}
	resp, err := getResponse(m)
	if err != nil {
		return redmineIssue, err
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return redmineIssue, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(byteArray, &redmineIssue); err != nil {
		log.Fatalf("Error!: %v", err)
	}
	return redmineIssue, err
}

func redmine(url string)(msg string) {
	res, _ := encodeJson4Redmine(url)
	fmt.Printf("ID: %d\n", res.Issue.ID)
	fmt.Println("Project: " + res.Issue.Project.Name)
	fmt.Println("Title: " + res.Issue.Subject)
	fmt.Println("Status: " + res.Issue.Status.Name)
	fmt.Println("Category: " + res.Issue.Category.Name)

	msg = "#" + strconv.Itoa(res.Issue.ID) + " (" + res.Issue.Status.Name + ") " + res.Issue.Project.Name + " : " + res.Issue.Subject + " [" + url +"]"

	return
}
