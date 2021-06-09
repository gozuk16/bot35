package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Keywords struct {
	Key string
}

type Redmine struct {
	Url      string
	APIToken string
	Keywords []Keywords
}

type Auth struct {
	Site  string
	Token string
}

type Exclude struct {
	Site string
}

type HttpSummary struct {
	Intra   string
	Auth    []Auth
	Exclude []Exclude
}

type Jira struct {
	Endpoint string
	Url      string
	APIUser  string
	APIToken string
	Keywords []Keywords
}

type Bitbucket struct {
	Endpoint string
	Url      string
	UserId   string
	Password string
	Keywords []Keywords
}

type Standard struct {
	Endpoint string
	Url      string
	Keywords []Keywords
}

type Config struct {
	BotId               string
	SlackAppToken       string
	SlackBotToken       string
	Redmine             Redmine
	HttpSummary         HttpSummary
	Jira                Jira
	Bitbucket           Bitbucket
	BitbucketPR         Standard
	QuestionsUnanswered Standard
	QuestionsList       Standard
}

var config Config

func run(api *slack.Client) int {

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
				msgs := []string{}

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

func setShuffle(txt string) string {
	m := strings.Split(strings.TrimSpace(txt), " ")[1:]
	log.Printf("m: %v\n", m)
	if len(m) == 1 {
		return "シャッフルですね。候補をスペース区切りで入れてください。 [ ex) shuffle a b c d ]"
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
		msg := "結果発表！\n"
		for i, v := range fortune {
			log.Printf("%v: %v\n", i, v)
			msg += " " + v
		}
		return msg
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

// setMessage
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

// createMessage キーワードを分解して渡された任意のfuncを呼んで結果を返す
func createMessage(t string, keywords []Keywords, endpoint string, fn func(string) (string, error)) string {
	pt, _, line, _ := runtime.Caller(0) // debug用に現在のスタックから情報を取得

	for _, k := range keywords {
		r := regexp.MustCompile(k.Key)
		str := r.FindAllStringSubmatch(t, -1)
		if str != nil {
			log.Printf("%s, %d: str len=%d\n", runtime.FuncForPC(pt).Name(), line, len(str))
			var m string
			for i, v := range str {
				log.Printf("%s, %d: str[%d]=%v\n", runtime.FuncForPC(pt).Name(), line, i, v[0])
				url := ""
				if len(v) > 1 {
					url = endpoint + v[1]
				} else {
					url = endpoint
				}
				if s, err := fn(url); err != nil {
					fmt.Println(err)
					m += err.Error()
				} else {
					m += s + "\n"
				}
			}
			return m
		}
	}
	return ""
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

func postMessage(api *slack.Client, event *slackevents.MessageEvent, msg string) {
	_, _, err := api.PostMessage(
		event.Channel,
		slack.MsgOptionText(msg, false),
	)
	if err != nil {
		log.Printf("Failed to reply: %v", err)
	}
}

func main() {
	getConfig()
	api := slack.New(
		config.SlackBotToken,
		slack.OptionAppLevelToken(config.SlackAppToken),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := api.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				// イベント API のハンドリング

				// 3 秒以内にとりあえず ack
				socketMode.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch ev := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if ev.User != selfUserId && (strings.Contains(ev.Text, "hello") ||
							strings.Contains(ev.Text, "こんにちは")) {
							msg := fmt.Sprintf(":wave: こんにちは <@%v> さん！", ev.User)
							postMessage(api, ev, msg)
						} else if ev.User != selfUserId && strings.Contains(ev.Text, "reload config") {
							postMessage(api, ev, "リロードするよ")
							getConfig()
						} else if ev.User != selfUserId && (strings.Contains(ev.Text, "おみくじ") ||
							strings.Contains(ev.Text, "shuffle") ||
							strings.Contains(ev.Text, "シャッフル") ||
							strings.Contains(ev.Text, "しゃっふる")) {
							msg := setShuffle(ev.Text)
							postMessage(api, ev, msg)
						} else if ev.User != selfUserId && strings.Contains(ev.Text, config.Redmine.Url) {
							// Redmine URL
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
							m := createMessage(text, config.Redmine.Keywords, config.Redmine.Url, redmine)
							postMessage(api, ev, m)

						} else if ev.User != selfUserId && (strings.Contains(ev.Text, "<http://") || strings.Contains(ev.Text, "<https://")) &&
							strings.Contains(ev.Text, config.HttpSummary.Intra) {
							var url string
							if !strings.Contains(ev.Text, config.Redmine.Url) {
								// http Summary
								key := "<(https?://.*." + config.HttpSummary.Intra + "/?.*?)>"
								log.Printf("key=%v\n", key)

								// exludeなら抜ける
								for _, ex := range config.HttpSummary.Exclude {
									if strings.Contains(ev.Text, ex) {
										break
									}
								}
								r := regexp.MustCompile(key)
								str := r.FindAllStringSubmatch(ev.Text, -1)
								log.Printf("str=%v\n", str)
								if str != nil {
									log.Printf("str len=%d\n", len(str))
									//var msg string
									for i, v := range str {
										log.Printf("str[%d]=%v\n", i, v[1])
										url = v[1]
										if url != "" {
											imgext := []string{"png", "jpg", "bmp"}
											for _, ext := range imgext {
												if strings.HasSuffix(url, ext) {
													log.Println("img: ", ext)
													if strings.HasPrefix(url, "https://jira.in.infocom.co.jp/redmine/attachments") {
														url = strings.Replace(url, "attachments", "attachments/download", 1)
														log.Println(url)
													}
													req, _ := http.NewRequest("GET", url, nil)
													req.Header.Set("X-Redmine-API-Key", config.Redmine.APIToken)
													client := new(http.Client)
													resp, err := client.Do(req)
													if err != nil {
														log.Println(err, url)
													} else {
														defer resp.Body.Close()

														_, err = api.UploadFile(
															slack.FileUploadParameters{
																Reader:   resp.Body,
																Filename: "image: " + url,
																Channels: []string{ev.Channel},
															})
														if err != nil {
															log.Println(err, url)
														}
													}
													url = "" // 画像の場合は空にして後でpostMessageしないように
												}
											}
										}
									}
									if url != "" {
										postMessage(api, ev, httpSummary(url))
									}
								}
							}
						} else if ev.User != selfUserId {
							// Redmine
							m := createMessage(ev.Text, config.Redmine.Keywords, config.Redmine.Url, redmine)
							if m != "" {
								postMessage(api, ev, m)
							}

							// JIRA
							m = createMessage(ev.Text, config.Jira.Keywords, config.Jira.Endpoint, jira)
							if m != "" {
								postMessage(api, ev, m)
							}

							// Bitbucket
							m = createMessage(ev.Text, config.Bitbucket.Keywords, config.Bitbucket.Endpoint, bitbucket)
							if m != "" {
								postMessage(api, ev, m)
							}
						}
					default:
						socketMode.Debugf("Skipped: %v", ev)
					}
				default:
					socketMode.Debugf("unsupported Events API eventPayload received")
				}
			}
		}
	}()

	socketMode.Run()
}
