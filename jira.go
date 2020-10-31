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
		Statuscategorychangedate string `json:"statuscategorychangedate"`
		Issuetype                struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			IconURL     string `json:"iconUrl"`
			Name        string `json:"name"`
			Subtask     bool   `json:"subtask"`
			AvatarID    int    `json:"avatarId"`
		} `json:"issuetype"`
		Timespent int `json:"timespent"`
		Project   struct {
			Self           string `json:"self"`
			ID             string `json:"id"`
			Key            string `json:"key"`
			Name           string `json:"name"`
			ProjectTypeKey string `json:"projectTypeKey"`
			Simplified     bool   `json:"simplified"`
			AvatarUrls     struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
		} `json:"project"`
		FixVersions []struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Archived    bool   `json:"archived"`
			Released    bool   `json:"released"`
			ReleaseDate string `json:"releaseDate"`
		} `json:"fixVersions"`
		Aggregatetimespent int `json:"aggregatetimespent"`
		Resolution         struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
		} `json:"resolution"`
		Customfield11401 interface{} `json:"customfield_11401"`
		Customfield11400 string      `json:"customfield_11400"`
		Customfield10500 interface{} `json:"customfield_10500"`
		Resolutiondate   string      `json:"resolutiondate"`
		Workratio        int         `json:"workratio"`
		Watches          struct {
			Self       string `json:"self"`
			WatchCount int    `json:"watchCount"`
			IsWatching bool   `json:"isWatching"`
		} `json:"watches"`
		Issuerestriction struct {
			Issuerestrictions struct {
			} `json:"issuerestrictions"`
			ShouldDisplay bool `json:"shouldDisplay"`
		} `json:"issuerestriction"`
		LastViewed string `json:"lastViewed"`
		Created    string `json:"created"`
		Priority   struct {
			Self    string `json:"self"`
			IconURL string `json:"iconUrl"`
			Name    string `json:"name"`
			ID      string `json:"id"`
		} `json:"priority"`
		Customfield10300              interface{}   `json:"customfield_10300"`
		Customfield10301              interface{}   `json:"customfield_10301"`
		Labels                        []string      `json:"labels"`
		Customfield11303              interface{}   `json:"customfield_11303"`
		Customfield11304              interface{}   `json:"customfield_11304"`
		Customfield11305              interface{}   `json:"customfield_11305"`
		Timeestimate                  int           `json:"timeestimate"`
		Aggregatetimeoriginalestimate int           `json:"aggregatetimeoriginalestimate"`
		Versions                      []interface{} `json:"versions"`
		Issuelinks                    []struct {
			ID   string `json:"id"`
			Self string `json:"self"`
			Type struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Inward  string `json:"inward"`
				Outward string `json:"outward"`
				Self    string `json:"self"`
			} `json:"type"`
			OutwardIssue struct {
				ID     string `json:"id"`
				Key    string `json:"key"`
				Self   string `json:"self"`
				Fields struct {
					Summary string `json:"summary"`
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
					Priority struct {
						Self    string `json:"self"`
						IconURL string `json:"iconUrl"`
						Name    string `json:"name"`
						ID      string `json:"id"`
					} `json:"priority"`
					Issuetype struct {
						Self        string `json:"self"`
						ID          string `json:"id"`
						Description string `json:"description"`
						IconURL     string `json:"iconUrl"`
						Name        string `json:"name"`
						Subtask     bool   `json:"subtask"`
						AvatarID    int    `json:"avatarId"`
					} `json:"issuetype"`
				} `json:"fields"`
			} `json:"outwardIssue"`
		} `json:"issuelinks"`
		Assignee struct {
			Self       string `json:"self"`
			AccountID  string `json:"accountId"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
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
		Description          interface{}   `json:"description"`
		Customfield11100     []string      `json:"customfield_11100"`
		Customfield11421     interface{}   `json:"customfield_11421"`
		Customfield11420     interface{}   `json:"customfield_11420"`
		Customfield11300     string        `json:"customfield_11300"`
		Customfield11301     string        `json:"customfield_11301"`
		Customfield11422     interface{}   `json:"customfield_11422"`
		Customfield11302     interface{}   `json:"customfield_11302"`
		Timetracking         struct {
			OriginalEstimate         string `json:"originalEstimate"`
			RemainingEstimate        string `json:"remainingEstimate"`
			TimeSpent                string `json:"timeSpent"`
			OriginalEstimateSeconds  int    `json:"originalEstimateSeconds"`
			RemainingEstimateSeconds int    `json:"remainingEstimateSeconds"`
			TimeSpentSeconds         int    `json:"timeSpentSeconds"`
		} `json:"timetracking"`
		Customfield10203 interface{} `json:"customfield_10203"`
		Customfield10005 string      `json:"customfield_10005"`
		Customfield11414 interface{} `json:"customfield_11414"`
		Customfield10600 string      `json:"customfield_10600"`
		Customfield11413 interface{} `json:"customfield_11413"`
		Customfield10402 struct {
			HasEpicLinkFieldDependency bool `json:"hasEpicLinkFieldDependency"`
			ShowField                  bool `json:"showField"`
			NonEditableReason          struct {
				Reason  string `json:"reason"`
				Message string `json:"message"`
			} `json:"nonEditableReason"`
		} `json:"customfield_10402"`
		Customfield10601      string        `json:"customfield_10601"`
		Security              interface{}   `json:"security"`
		Customfield11415      string        `json:"customfield_11415"`
		Attachment            []interface{} `json:"attachment"`
		Aggregatetimeestimate int           `json:"aggregatetimeestimate"`
		Customfield11418      interface{}   `json:"customfield_11418"`
		Customfield11419      interface{}   `json:"customfield_11419"`
		Summary               string        `json:"summary"`
		Creator               struct {
			Self       string `json:"self"`
			AccountID  string `json:"accountId"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
		} `json:"creator"`
		Subtasks []interface{} `json:"subtasks"`
		Reporter struct {
			Self       string `json:"self"`
			AccountID  string `json:"accountId"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
		} `json:"reporter"`
		Customfield10000  interface{} `json:"customfield_10000"`
		Aggregateprogress struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
			Percent  int `json:"percent"`
		} `json:"aggregateprogress"`
		Customfield11410 interface{} `json:"customfield_11410"`
		Customfield10200 interface{} `json:"customfield_10200"`
		Customfield11412 interface{} `json:"customfield_11412"`
		Customfield10201 interface{} `json:"customfield_10201"`
		Customfield10004 interface{} `json:"customfield_10004"`
		Customfield10202 interface{} `json:"customfield_10202"`
		Customfield11411 interface{} `json:"customfield_11411"`
		Customfield11403 interface{} `json:"customfield_11403"`
		Customfield11402 interface{} `json:"customfield_11402"`
		Environment      interface{} `json:"environment"`
		Customfield11405 interface{} `json:"customfield_11405"`
		Customfield11404 interface{} `json:"customfield_11404"`
		Customfield11407 interface{} `json:"customfield_11407"`
		Customfield11406 interface{} `json:"customfield_11406"`
		Customfield11409 interface{} `json:"customfield_11409"`
		Duedate          string      `json:"duedate"`
		Customfield11408 interface{} `json:"customfield_11408"`
		Progress         struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
			Percent  int `json:"percent"`
		} `json:"progress"`
		Votes struct {
			Self     string `json:"self"`
			Votes    int    `json:"votes"`
			HasVoted bool   `json:"hasVoted"`
		} `json:"votes"`
		Comment struct {
			Comments   []interface{} `json:"comments"`
			MaxResults int           `json:"maxResults"`
			Total      int           `json:"total"`
			StartAt    int           `json:"startAt"`
		} `json:"comment"`
		Worklog struct {
			StartAt    int `json:"startAt"`
			MaxResults int `json:"maxResults"`
			Total      int `json:"total"`
			Worklogs   []struct {
				Self   string `json:"self"`
				Author struct {
					Self       string `json:"self"`
					AccountID  string `json:"accountId"`
					AvatarUrls struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
					AccountType string `json:"accountType"`
				} `json:"author"`
				UpdateAuthor struct {
					Self       string `json:"self"`
					AccountID  string `json:"accountId"`
					AvatarUrls struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
					AccountType string `json:"accountType"`
				} `json:"updateAuthor"`
				Comment struct {
					Version int    `json:"version"`
					Type    string `json:"type"`
					Content []struct {
						Type    string `json:"type"`
						Content []struct {
							Type  string `json:"type"`
							Text  string `json:"text,omitempty"`
							Attrs struct {
								URL string `json:"url"`
							} `json:"attrs,omitempty"`
						} `json:"content"`
					} `json:"content"`
				} `json:"comment"`
				Created          string `json:"created"`
				Updated          string `json:"updated"`
				Started          string `json:"started"`
				TimeSpent        string `json:"timeSpent"`
				TimeSpentSeconds int    `json:"timeSpentSeconds"`
				ID               string `json:"id"`
				IssueID          string `json:"issueId"`
			} `json:"worklogs"`
		} `json:"worklog"`
	} `json:"fields"`
}

var jiraIssue JiraIssue

func encodeJson4Jira(url string) (JiraIssue, error) {
	m := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"url":          url}
	resp, err := getResponseWithBasicAuth(m, config.Jira.APIUser, config.Jira.APIToken)
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

func jira(url string) (msg string, err error) {
	res, err := encodeJson4Jira(url)
	fmt.Println("Key: " + res.Key)
	fmt.Println("Project: " + res.Fields.Project.Name)
	fmt.Println("Title: " + res.Fields.Summary)
	fmt.Println("Status: " + res.Fields.Status.Name)
	fmt.Println("Issuetype: " + res.Fields.Issuetype.Name)

	ticketUrl := config.Jira.Url + res.Key
	msg = res.Key + " (" + res.Fields.Status.Name + ") " + res.Fields.Project.Name + " : " + res.Fields.Summary + " [" + ticketUrl + "]"

	return msg, err
}
