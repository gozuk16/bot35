build: bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	go build bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go

linux: bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	GOOS=linux GOARCH=amd64 go build -o linux-amd64/bot35

win: bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	GOOS=windows GOARCH=amd64 go build -o windows-amd64/bot35.exe

run: bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
	go run bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go
