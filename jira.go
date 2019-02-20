package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type JiraIssue struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Issuetype struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			IconURL     string `json:"iconUrl"`
			Name        string `json:"name"`
			Subtask     bool   `json:"subtask"`
			AvatarID    int    `json:"avatarId"`
		} `json:"issuetype"`
		Timespent interface{} `json:"timespent"`
		Project   struct {
			Self       string `json:"self"`
			ID         string `json:"id"`
			Key        string `json:"key"`
			Name       string `json:"name"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			ProjectCategory struct {
				Self        string `json:"self"`
				ID          string `json:"id"`
				Description string `json:"description"`
				Name        string `json:"name"`
			} `json:"projectCategory"`
		} `json:"project"`
		Customfield11000   interface{}   `json:"customfield_11000"`
		FixVersions        []interface{} `json:"fixVersions"`
		Customfield11001   interface{}   `json:"customfield_11001"`
		Customfield11002   interface{}   `json:"customfield_11002"`
		Aggregatetimespent interface{}   `json:"aggregatetimespent"`
		Resolution         interface{}   `json:"resolution"`
		Customfield11003   interface{}   `json:"customfield_11003"`
		Customfield10500   interface{}   `json:"customfield_10500"`
		Customfield10700   interface{}   `json:"customfield_10700"`
		Customfield10701   []struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
		} `json:"customfield_10701"`
		Customfield10900 interface{} `json:"customfield_10900"`
		Resolutiondate   interface{} `json:"resolutiondate"`
		Workratio        int         `json:"workratio"`
		LastViewed       interface{} `json:"lastViewed"`
		Watches          struct {
			Self       string `json:"self"`
			WatchCount int    `json:"watchCount"`
			IsWatching bool   `json:"isWatching"`
		} `json:"watches"`
		Created  string `json:"created"`
		Priority struct {
			Self    string `json:"self"`
			IconURL string `json:"iconUrl"`
			Name    string `json:"name"`
			ID      string `json:"id"`
		} `json:"priority"`
		Customfield10300              interface{}   `json:"customfield_10300"`
		Labels                        []interface{} `json:"labels"`
		Customfield10301              interface{}   `json:"customfield_10301"`
		Timeestimate                  int           `json:"timeestimate"`
		Aggregatetimeoriginalestimate int           `json:"aggregatetimeoriginalestimate"`
		Versions                      []interface{} `json:"versions"`
		Issuelinks                    []interface{} `json:"issuelinks"`
		Assignee                      struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
		} `json:"assignee"`
		Updated string `json:"updated"`
		Status  struct {
			Self           string `json:"self"`
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				Self      string `json:"self"`
				ID        int    `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
		Components           []interface{} `json:"components"`
		Timeoriginalestimate int           `json:"timeoriginalestimate"`
		Description          string        `json:"description"`
		Timetracking         struct {
			OriginalEstimate         string `json:"originalEstimate"`
			RemainingEstimate        string `json:"remainingEstimate"`
			OriginalEstimateSeconds  int    `json:"originalEstimateSeconds"`
			RemainingEstimateSeconds int    `json:"remainingEstimateSeconds"`
		} `json:"timetracking"`
		Customfield10203      interface{}   `json:"customfield_10203"`
		Customfield10005      string        `json:"customfield_10005"`
		Customfield10402      interface{}   `json:"customfield_10402"`
		Customfield10600      interface{}   `json:"customfield_10600"`
		Customfield10601      interface{}   `json:"customfield_10601"`
		Customfield10800      string        `json:"customfield_10800"`
		Attachment            []interface{} `json:"attachment"`
		Aggregatetimeestimate int           `json:"aggregatetimeestimate"`
		Summary               string        `json:"summary"`
		Creator               struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
		} `json:"creator"`
		Subtasks []interface{} `json:"subtasks"`
		Reporter struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
		} `json:"reporter"`
		Customfield10000  interface{} `json:"customfield_10000"`
		Aggregateprogress struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
			Percent  int `json:"percent"`
		} `json:"aggregateprogress"`
		Customfield10200 string      `json:"customfield_10200"`
		Customfield10201 interface{} `json:"customfield_10201"`
		Customfield10202 interface{} `json:"customfield_10202"`
		Customfield10004 []interface{} `json:"customfield_10004"`
		Environment      interface{} `json:"environment"`
		Duedate          interface{} `json:"duedate"`
		Progress         struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
			Percent  int `json:"percent"`
		} `json:"progress"`
		Comment struct {
			Comments []struct {
				Self   string `json:"self"`
				ID     string `json:"id"`
				Author struct {
					Self         string `json:"self"`
					Name         string `json:"name"`
					Key          string `json:"key"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls   struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
				} `json:"author"`
				Body         string `json:"body"`
				UpdateAuthor struct {
					Self         string `json:"self"`
					Name         string `json:"name"`
					Key          string `json:"key"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls   struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
				} `json:"updateAuthor"`
				Created string `json:"created"`
				Updated string `json:"updated"`
			} `json:"comments"`
			MaxResults int `json:"maxResults"`
			Total      int `json:"total"`
			StartAt    int `json:"startAt"`
		} `json:"comment"`
		Votes struct {
			Self     string `json:"self"`
			Votes    int    `json:"votes"`
			HasVoted bool   `json:"hasVoted"`
		} `json:"votes"`
		Worklog struct {
			StartAt    int           `json:"startAt"`
			MaxResults int           `json:"maxResults"`
			Total      int           `json:"total"`
			Worklogs   []interface{} `json:"worklogs"`
		} `json:"worklog"`
	} `json:"fields"`
}

var jiraIssue JiraIssue

func encodeJson4Jira(url string) (JiraIssue, error) {
	m := map[string]string{"url": url}
	resp, err := getResponse(m)
	if err != nil {
		return jiraIssue, err
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jiraIssue, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(byteArray, &jiraIssue); err != nil {
		log.Fatalf("Error!: %v", err)
	}
	return jiraIssue, err
}

func jira(url string) (msg string) {
	res, _ := encodeJson4Jira(url)
	fmt.Println("Key: " + res.Key)
	fmt.Println("Project: " + res.Fields.Project.Name)
	fmt.Println("Title: " + res.Fields.Summary)
	fmt.Println("Status: " + res.Fields.Status.Name)
	fmt.Println("Issuetype: " + res.Fields.Issuetype.Name)

	ticketUrl := config.Jira.Url + res.Key
	msg = res.Key + " (" + res.Fields.Status.Name + ") " + res.Fields.Project.Name + " : " + res.Fields.Summary + " [" + ticketUrl + "]"

	return
}
