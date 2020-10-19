SRC=bot35.go httpclient.go redmine.go jira.go httpSummary.go bitbucket.go bitbucket_pr.go confluence.go

build:
	go build $(SRC)

linux:
	GOOS=linux GOARCH=amd64 go build -o linux-amd64/bot35 $(SRC)

win:
	GOOS=windows GOARCH=amd64 go build -o windows-amd64/bot35.exe $(SRC)

run:
	go run $(SRC)
