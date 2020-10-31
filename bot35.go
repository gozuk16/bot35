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
	APIUser  string
	APIToken string
	Keywords []Keywords
}

/*
type Bitbucket struct {
	Endpoint string
	Url      string
	Keywords []Keywords
}
*/

type Standard struct {
	Endpoint string
	Url      string
	Keywords []Keywords
}

type Config struct {
	BotId               string
	SlackAPIToken       string
	Redmine             Redmine
	HttpSummary         HttpSummary
	Jira                Jira
	Bitbucket           Standard
	BitbucketPR         Standard
	QuestionsUnanswered Standard
	QuestionsList       Standard
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

				msgs := []string{}

				if ev.Text == "こんにちは" || ev.Text == "hello" {
					msgs = append(msgs, user.Profile.RealName+"さん、こんにちは(^-^)")
				}

				if ev.Text == "ぬるぽ" || ev.Text == "NullPointerException" {
					msgs = append(msgs, nullpo)
				}

				// Redmine
				if strings.Contains(ev.Text, config.Redmine.Url) {
					log.Println(ev.Text)
					var text string
					r := regexp.MustCompile(config.Redmine.Url + "([0-9]*)")
					str := r.FindAllStringSubmatch(ev.Text, -1)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						for i, v := range str {
							log.Printf("str[%d]0=%v\n", i, v[0])
							log.Printf("str[%d]1=%v\n", i, v[1])
							text = "redmine " + v[1]
							break
						}
					}
					setMessage(text, config.Redmine.Keywords, config.Redmine.Url, redmine, &msgs)
				} else {
					setMessage(ev.Text, config.Redmine.Keywords, config.Redmine.Url, redmine, &msgs)
				}

				// JIRA
				setMessage(ev.Text, config.Jira.Keywords, config.Jira.Endpoint, jira, &msgs)

				// Bitbucket
				setMessage(ev.Text, config.Bitbucket.Keywords, config.Bitbucket.Endpoint, bitbucket, &msgs)

				// Bitbucket Pull-Request
				for _, k := range config.BitbucketPR.Keywords {
					r := regexp.MustCompile(k.Key)
					str := r.FindAllStringSubmatch(ev.Text, -1)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						var msg string
						for i, v := range str {
							log.Printf("str[%d]0=%v\n", i, v[0])
							log.Printf("str[%d]1=%v\n", i, v[1])
							log.Printf("str[%d]2=%v\n", i, v[2])
							str := strings.Replace(config.BitbucketPR.Endpoint, "{0}", v[1], 1)
							bitbucketPRApi := strings.Replace(str, "{1}", v[2], 1)
							msg += pr(bitbucketPRApi) + "\n"
						}
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}
				}

				// Questions for Confluence
				//setMessage(ev.Text, config.QuestionsUnanswered.Keywords, config.QuestionsUnanswered.Endpoint, confluence, &msgs)
				for _, k := range config.QuestionsUnanswered.Keywords {
					r := regexp.MustCompile(k.Key)
					str := r.FindAllStringSubmatch(ev.Text, -1)
					if str != nil {
						log.Printf("str len=%d\n", len(str))
						var msg string
						for i, v := range str {
							log.Printf("str[%d]=%v\n", i, v[0])
							//confluenceApi := config.Confluence.Endpoint + v[1]
							confluenceApi := config.QuestionsUnanswered.Endpoint
							msg += confluence(confluenceApi) + "\n"
						}
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}
				}

				// http Summary
				if (strings.Contains(ev.Text, "<http://") || strings.Contains(ev.Text, "<https://")) &&
					(strings.Contains(ev.Text, config.HttpSummary.Intra) && strings.Contains(ev.Text, config.Redmine.Url) == false) {

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

				// メンション対応
				responseMention(botId, ev.Text, &msgs)

				for _, m := range msgs {
					rtm.SendMessage(rtm.NewOutgoingMessage(m, ev.Channel))
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func responseMention(botId string, txt string, msgs *[]string) {
	log.Printf("botId: %v\n", botId)
	if strings.HasPrefix(txt, "<@"+botId+">") {
		m := strings.Split(strings.TrimSpace(txt), " ")[1:]
		log.Printf("m: %v\n", m)
		if len(m) == 0 {
			*msgs = append(*msgs, "呼んだ？")
		} else if m[0] == "reload" {
			*msgs = append(*msgs, "りろーどするよ")
			getConfig()
		} else if m[0] == "おみくじ" || m[0] == "shuffle" || m[0] == "シャッフル" || m[0] == "しゃっふる" {
			if len(m) == 1 {
				*msgs = append(*msgs, "シャッフルですね。候補をスペース区切りで入れてください。 [ @botsan shuffle a b c d ]")
			} else {
				var fortune []string
				for i, v := range m {
					if i == 0 {
						continue
					}
					log.Printf("%v: %v\n", i, v)
					fortune = append(fortune, v)
				}
				shuffle(fortune)
				log.Printf("shuffle\n")
				*msgs = append(*msgs, "結果発表！")
				for i, v := range fortune {
					log.Printf("%v: %v\n", i, v)
					*msgs = append(*msgs, v)
				}
			}
		}
	}

}

func setMessage(txt string, keywords []Keywords, endpoint string, fn func(string) (string, error), msgs *[]string) {
	for _, k := range keywords {
		r := regexp.MustCompile(k.Key)
		str := r.FindAllStringSubmatch(txt, -1)
		if str != nil {
			log.Printf("str len=%d\n", len(str))
			var m string
			for i, v := range str {
				log.Printf("str[%d]=%v\n", i, v[0])
				url := ""
				if len(v) > 1 {
					url = endpoint + v[1]
				} else {
					url = endpoint
				}
				if s, err := fn(url); err != nil {
					log.Println(err)
					m += err.Error()
				} else {
					m += s + "\n"
				}
			}
			*msgs = append(*msgs, m)
		}
	}
}

func shuffle(data []string) {
	n := len(data)
	log.Printf("n=%v\n", n)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		log.Printf("j=%v, i+1=%v", j, i+1)
		data[i], data[j] = data[j], data[i]
	}
}

func getConfig() {
	configFile := "bot35.json"
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
