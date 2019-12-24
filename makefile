NO_PROXY=foo.samples.com

build: botsango.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	go build botsango.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
run: botsango.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	go run botsango.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
