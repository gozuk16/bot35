package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/pkg/errors"
)

/*
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
*/

type RedmineIssueWithJournals struct {
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
		FixedVersion struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"fixed_version"`
		Subject     string      `json:"subject"`
		Description string      `json:"description"`
		StartDate   interface{} `json:"start_date"`
		DueDate     interface{} `json:"due_date"`
		DoneRatio   int         `json:"done_ratio"`
		//IsPrivate           bool        `json:"is_private"`
		//EstimatedHours      interface{} `json:"estimated_hours"`
		//TotalEstimatedHours interface{} `json:"total_estimated_hours"`
		CustomFields []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"custom_fields"`
		CreatedOn time.Time `json:"created_on"`
		UpdatedOn time.Time `json:"updated_on"`
		ClosedOn  time.Time `json:"closed_on"`
		Journals  []struct {
			ID   int `json:"id"`
			User struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
			Notes     string    `json:"notes"`
			CreatedOn time.Time `json:"created_on"`
			//PrivateNotes bool      `json:"private_notes"`
			/*
				Details      []struct {
					Property string `json:"property"`
					Name     string `json:"name"`
					OldValue string `json:"old_value"`
					NewValue string `json:"new_value"`
				} `json:"details"`
			*/
		} `json:"journals"`
	} `json:"issue"`
}

var redmineIssueWithJournals RedmineIssueWithJournals

func encodeJson4Redmine(uri string) (RedmineIssueWithJournals, error) {
	m := map[string]string{
		"url":         uri,
		"headerKey":   "X-Redmine-API-Key",
		"headerValue": config.Redmine.APIToken}
	resp, err := getResponse(m)
	if err != nil {
		return redmineIssueWithJournals, err
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return redmineIssueWithJournals, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(byteArray, &redmineIssueWithJournals); err != nil {
		log.Fatalf("Error!: %v", err)
	}
	fmt.Printf("%v\n", redmineIssueWithJournals)
	return redmineIssueWithJournals, err
}

func redmine(uri string) (msg string, err error) {
	res, err := encodeJson4Redmine(uri + config.Redmine.EndpointParam)
	log.Println("url:", uri)
	log.Println("api url:", uri+config.Redmine.EndpointParam)
	log.Printf("ID: %d\n", res.Issue.ID)
	log.Println("Project:", res.Issue.Project.Name)
	log.Println("Title:", res.Issue.Subject)
	log.Println("Status:", res.Issue.Status.Name)
	log.Println("Category:", res.Issue.Category.Name)

	msg = fmt.Sprintf("#%d (%s):%s%s\n[%s]",
		res.Issue.ID,
		res.Issue.Status.Name,
		res.Issue.Project.Name,
		res.Issue.Subject,
		uri)

	return msg, err
}

func redmineNote(uri string, no int) (msg string, err error) {
	res, err := encodeJson4Redmine(uri + config.Redmine.EndpointParam)
	if err != nil {
		return "", err
	}

	if len(res.Issue.Journals) < no {
		return "", errors.Errorf("journals index error. journals len: %d, no: %d", len(res.Issue.Journals), no)
	}

	log.Println("url:", uri)
	log.Println("api url:", uri+config.Redmine.EndpointParam)
	log.Printf("ID: %d\n", res.Issue.ID)
	log.Println("Project:", res.Issue.Project.Name)
	log.Println("Title:", res.Issue.Subject)
	log.Println("Status:", res.Issue.Status.Name)
	log.Println("Category:", res.Issue.Category.Name)

	sort.SliceStable(res.Issue.Journals, func(i, j int) bool { return res.Issue.Journals[i].ID < res.Issue.Journals[j].ID })
	for i, v := range res.Issue.Journals {
		log.Printf("Journals[%d] %s: %s\n", i, v.User.Name, v.Notes)
	}

	loc, _ := time.LoadLocation("Local")
	jst := res.Issue.Journals[no].CreatedOn.In(loc)
	log.Printf("CreatedOn UTC: %v\n", res.Issue.Journals[no].CreatedOn)
	log.Printf("CreatedOn JST: %v\n", jst)

	msg = fmt.Sprintf("#%d (%s):%s%s\n--- %s さんが%vに更新 ---\n%s",
		res.Issue.ID,
		res.Issue.Status.Name,
		res.Issue.Project.Name,
		res.Issue.Subject,
		res.Issue.Journals[no].User.Name,
		jst,
		res.Issue.Journals[no].Notes)

	return msg, err
}
