package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type BitbucketRepos struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Values     []struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Origin        struct {
			Slug          string `json:"slug"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			ScmID         string `json:"scmId"`
			State         string `json:"state"`
			StatusMessage string `json:"statusMessage"`
			Forkable      bool   `json:"forkable"`
			Project       struct {
				Key         string `json:"key"`
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Public      bool   `json:"public"`
				Type        string `json:"type"`
				Links       struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"project"`
			Public bool `json:"public"`
			Links  struct {
				Clone []struct {
					Href string `json:"href"`
					Name string `json:"name"`
				} `json:"clone"`
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"origin"`
		Project struct {
			Key   string `json:"key"`
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Owner struct {
				Name         string `json:"name"`
				EmailAddress string `json:"emailAddress"`
				ID           int    `json:"id"`
				DisplayName  string `json:"displayName"`
				Active       bool   `json:"active"`
				Slug         string `json:"slug"`
				Type         string `json:"type"`
				Links        struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"owner"`
			Links struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public bool `json:"public"`
		Links  struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"values"`
	Start int `json:"start"`
}

var bitbucketRepos BitbucketRepos

func encodeJson4Bitbucket(url string) (BitbucketRepos, error) {
	m := map[string]string{"url": url}
	resp, err := getResponse(m)
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
	fmt.Printf("Num: %d\n", res.Size)
	msg = ""
	for _, v := range res.Values {
		//fmt.Printf("Project: %v\n", res.Values[i].Origin.Project.Key)
		fmt.Printf("v.Project: %v\n", v.Project)
		msg += v.Project.Name + "/"
		fmt.Printf("v.Name: %v\n", v.Name)
		msg += v.Name + " "
		fmt.Printf("v.Links.Self[0]: %v\n", v.Links.Self[0].Href)
		msg += "[" + v.Links.Self[0].Href + "]\n"
	}

	//msg = ""
	//msg = res.Key + " (" + res.Fields.Status.Name + ") " + res.Fields.Project.Name + " : " + res.Fields.Summary + " [" + ticketUrl + "]"

	return msg, err
}
