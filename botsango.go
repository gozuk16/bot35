package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type Keywords struct {
	Key string
}

type Redmine struct {
	Url      string
	APIToken string
	Keywords []Keywords
}

type HttpSummary struct {
	Intra string
}

type Jira struct {
	Endpoint string
	Url      string
	Keywords []Keywords
}

type Bitbucket struct {
	Url      string
	Keywords []Keywords
}

type Confluence struct {
	Url      string
	Keywords []Keywords
}

type Config struct {
	BotId         string
	SlackAPIToken string
	Redmine       Redmine
	HttpSummary   HttpSummary
	Jira          Jira
	Bitbucket     Bitbucket
	Confluence    Confluence
}

var config Config

func run(api *slack.Client) int {

	nullpo := "```"
	nullpo += `
　　 （　・∀・）　　　|　|　ｶﾞｯ
　　と　　　　）　 　 |　|
　　　 Ｙ　/ノ　　　 人
　　　　 /　）　 　 < 　>__Λ∩
　　 ＿/し'　／／. Ｖ｀Д´）/ ←>>1
　　（＿フ彡　　　　　 　　/
`
	nullpo += "```"

	//botId := "U1Y4HGEJU"
	botId := config.BotId

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				log.Printf("ev.Info.User.ID: %v\n", ev.Info.User.ID)
				log.Printf("ev.Info.User.Name: %v\n", ev.Info.User.Name)

			case *slack.HelloEvent:
				log.Print("Hello Event")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)
				log.Printf("ev.User: %v\n", ev.User)
				log.Printf("ev.Text: %v\n", ev.Text)
				log.Printf("ev.Channel: %v\n", ev.Channel)
				log.Printf("ev.Msg.Text: %v\n", ev.Msg.Text)

				user, _ := api.GetUserInfo(ev.User)
				log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
				if strings.HasPrefix(ev.Text, "こんにちは") {
					rtm.SendMessage(rtm.NewOutgoingMessage(user.Profile.RealName+"さん、こんにちは(^-^)", ev.Channel))
				}
				if ev.Text == "ぬるぽ" || ev.Text == "NullPointerException" {
					rtm.SendMessage(rtm.NewOutgoingMessage(nullpo, ev.Channel))
				}

				// Redmine
				for _, k := range config.Redmine.Keywords {
					r := regexp.MustCompile(k.Key)
					str := r.FindAllStringSubmatch(ev.Text, -1)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						var msg string
						for i, v := range str {
							log.Printf("str[%d]=%v\n", i, v[0])
							redmineUrl := config.Redmine.Url + v[1]
							msg += redmine(redmineUrl) + "\n"
						}
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}

				}

				// JIRA
				for _, k := range config.Jira.Keywords {
					r := regexp.MustCompile(k.Key)
					str := r.FindAllStringSubmatch(ev.Text, -1)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						var msg string
						for i, v := range str {
							log.Printf("str[%d]=%v\n", i, v[0])
							jiraApi := config.Jira.Endpoint + v[1]
							msg += jira(jiraApi) + "\n"
						}
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}

				}

				// HTTP Summary
				if strings.Contains(ev.Text, "<http://") || strings.Contains(ev.Text, "<https://") {
					key := "<(https?://.*." + config.HttpSummary.Intra + "/?.*?)>"
					log.Printf("key=%v\n", key)
					r := regexp.MustCompile(key)
					str := r.FindAllStringSubmatch(ev.Text, -1)
					log.Printf("str=%v\n", str)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						//var msg string
						for i, v := range str {
							log.Printf("str[%d]=%v\n", i, v[1])
							rtm.SendMessage(rtm.NewOutgoingMessage(httpSummary(v[1]), ev.Channel))
						}
					}
				}

				// おみくじ
				if strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", botId)) {
					m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
					log.Printf("m: %v\n", m)
					if len(m) == 0 {
						log.Printf("invalid message")
					} else {
						var fortune []string
						for i, v := range m[1:] {
							log.Printf("%v: %v\n", i, v)
							fortune = append(fortune, v)
						}
						shuffle(fortune)
						log.Printf("shuffle\n")
						for i, v := range fortune {
							log.Printf("%v: %v\n", i, v)
						}
						rtm.SendMessage(rtm.NewOutgoingMessage("おみくじですね。候補者を入れてください。", ev.Channel))
					}
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func shuffle(data []string) {
	n := len(data)
	log.Printf("n=%v", n)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		log.Printf("j=%v, i+1=%v", j, i+1)
		data[i], data[j] = data[j], data[i]
	}
}

func getConfig() {
	configFile := "botsango.json"
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		fmt.Println(p)
		f := filepath.Join(p, configFile)
		_, err := os.Stat(f)
		if err == nil {
			fmt.Println(f)
			configFile = f
			break
		}
	}
	jsonBuffer, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("ERROR: ", err)
		os.Exit(1)
	}

	err = json.Unmarshal(jsonBuffer, &config)
	if err != nil {
		log.Println("ERROR: ", err)
		os.Exit(1)
	}

	log.Printf("%#v\n", config)
}

func main() {
	getConfig()
	log.Printf("config.SlackAPIToken: %v\n", config.SlackAPIToken)
	api := slack.New(config.SlackAPIToken)
	//api.SetDebug(true)
	os.Exit(run(api))
}
