package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type BitbucketRepos struct {
	Scm     string      `json:"scm"`
	Website interface{} `json:"website"`
	HasWiki bool        `json:"has_wiki"`
	UUID    string      `json:"uuid"`
	Links   struct {
		Watchers struct {
			Href string `json:"href"`
		} `json:"watchers"`
		Branches struct {
			Href string `json:"href"`
		} `json:"branches"`
		Tags struct {
			Href string `json:"href"`
		} `json:"tags"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Source struct {
			Href string `json:"href"`
		} `json:"source"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		Hooks struct {
			Href string `json:"href"`
		} `json:"hooks"`
		Forks struct {
			Href string `json:"href"`
		} `json:"forks"`
		Downloads struct {
			Href string `json:"href"`
		} `json:"downloads"`
		Pullrequests struct {
			Href string `json:"href"`
		} `json:"pullrequests"`
	} `json:"links"`
	ForkPolicy string `json:"fork_policy"`
	FullName   string `json:"full_name"`
	Name       string `json:"name"`
	Project    struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		Type string `json:"type"`
		Name string `json:"name"`
		Key  string `json:"key"`
		UUID string `json:"uuid"`
	} `json:"project"`
	Language   string    `json:"language"`
	CreatedOn  time.Time `json:"created_on"`
	Mainbranch struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"mainbranch"`
	Workspace struct {
		Slug  string `json:"slug"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		UUID string `json:"uuid"`
	} `json:"workspace"`
	HasIssues bool `json:"has_issues"`
	Owner     struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	} `json:"owner"`
	UpdatedOn   time.Time `json:"updated_on"`
	Size        int       `json:"size"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	IsPrivate   bool      `json:"is_private"`
	Description string    `json:"description"`
}

var bitbucketRepos BitbucketRepos

func encodeJson4Bitbucket(url string) (BitbucketRepos, error) {
	m := map[string]string{"url": url}
	//resp, err := getResponse(m)
	resp, err := getResponseWithBasicAuth(m, config.Bitbucket.UserId, config.Bitbucket.Password)
	if err != nil {
		return bitbucketRepos, err
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bitbucketRepos, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(byteArray, &bitbucketRepos); err != nil {
		log.Fatalf("Error!: %v", err)
	}
	return bitbucketRepos, err
}

func bitbucket(url string) (msg string, err error) {
	res, err := encodeJson4Bitbucket(url)
	fmt.Printf("Scm: %s\n", res.Scm)
	fmt.Printf("Slug: %s\n", res.Slug)

	fmt.Printf("Project: %v\n", res.Project)
	msg += res.Project.Name + "/"
	fmt.Printf("Name: %v\n", res.Name)
	msg += res.Name + " "
	fmt.Printf("Links.Self: %v\n", res.Links)
	msg += "[" + res.Links.HTML.Href + "]\n"

	return msg, err
}
