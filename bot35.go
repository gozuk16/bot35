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
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Keywords struct {
	Key string
}

type Redmine struct {
	EndpointParam string
	Url           string
	APIToken      string
	Keywords      []Keywords
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

func isShuffle(selfUserId string, user string, text string) bool {
	if user != selfUserId && (strings.Contains(text, "おみくじ") ||
		strings.Contains(text, "shuffle") ||
		strings.Contains(text, "シャッフル") ||
		strings.Contains(text, "しゃっふる")) {
		return true
	}
	return false
}

// setShuffle ランダムに並べ替えた順番を返す
func setShuffle(txt string) string {
	// spaceと→の両方を区切り文字とする
	reg := "[ →]"

	// 1つ目はコマンドなんで取り去る
	m := regexp.MustCompile(reg).Split(txt, -1)[1:]
	log.Printf("text split: len=%d, text=%v\n", len(m), m)
	if len(m) <= 1 {
		//return fmt.Sprintf("%sをシャッフルできません。右のように入力してください。 [shuffle a b c d | shuffle a→b→c→d]", txt)
		return fmt.Sprintf("シャッフルできません。右のように入力してください。\n%s -> [shuffle a b c d | shuffle a→b→c→d]", txt)
	} else {
		var member []string
		for i, v := range m {
			log.Printf("%v: %v\n", i, v)
			member = append(member, v)
		}
		shuffle(member)
		log.Printf("shuffle\n")
		msg := "結果発表！\n"
		for i, v := range member {
			log.Printf("%v: %v\n", i, v)
			if i == 0 {
				msg += v
			} else {
				msg += "→" + v
			}
		}
		return msg
	}
}

// shuffle Fisher-Yates shuffle
func shuffle(data []string) {
	n := len(data)
	log.Printf("n=%v\n", n)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		log.Printf("j=%v, i+1=%v", j, i+1)
		data[i], data[j] = data[j], data[i]
	}
}

// isWebPage textにWebページが含まれているか判定する
func isWebPage(text string) bool {
	if strings.Contains(text, "<http://") || strings.Contains(text, "<https://") {
		return true
	}
	return false
}

// isIntraWebPage textにイントラのドメインが含まれているか判定する
func isIntraWebPage(text string) bool {
	if strings.Contains(text, config.HttpSummary.Intra) {
		return true
	}
	return false
}

func getImage(api *slack.Client, ev *slackevents.MessageEvent, url *string) {
	log.Println(*url)
	if *url == "" {
		return
	}

	imgext := []string{"png", "jpg", "bmp"}
	for _, ext := range imgext {
		if strings.HasSuffix(strings.ToLower(*url), ext) {
			log.Println("img: ", ext)
			req, _ := http.NewRequest("GET", *url, nil)
			// ToDo: 切り出して設定化
			if strings.Contains(*url, "https://jira.in.infocom.co.jp/redmine/attachments") {
				//*url = strings.Replace(*url, "attachments", "attachments/download", 1)
				req.Header.Set("X-Redmine-API-Key", config.Redmine.APIToken)
			}
			client := new(http.Client)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err, *url)
			} else {
				defer resp.Body.Close()

				_, err = api.UploadFile(
					slack.FileUploadParameters{
						Reader:   resp.Body,
						Filename: "image: " + *url,
						Channels: []string{ev.Channel},
					})
				if err != nil {
					log.Println(err, *url)
				}
			}
			*url = "" // 画像の場合は空にして後でpostMessageしないように
		}
	}
}

func isHello(selfUserId string, user string, text string) bool {
	if user != selfUserId && (strings.Contains(text, "hello") || strings.Contains(text, "こんにちは")) {
		return true
	}
	return false
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

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("ERROR:", err)
		os.Exit(1)
	}

	f := filepath.Join(pwd, configFile)
	_, err = os.Stat(f)
	if err == nil {
		log.Println(f)
		configFile = f
	}
	jsonBuffer, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("ERROR:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(jsonBuffer, &config)
	if err != nil {
		log.Println("ERROR:", err)
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
	setLogger()
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
						if isHello(selfUserId, ev.User, ev.Text) {
							msg := fmt.Sprintf(":wave: こんにちは <@%v> さん！", ev.User)
							postMessage(api, ev, msg)
						} else if ev.User != selfUserId && strings.Contains(ev.Text, "reload config") {
							postMessage(api, ev, "リロードするよ")
							getConfig()
						} else if isShuffle(selfUserId, ev.User, ev.Text) {
							msg := setShuffle(ev.Text)
							postMessage(api, ev, msg)
						} else if ev.User != selfUserId && strings.Contains(ev.Text, config.Redmine.Url) {
							// Redmine URL
							var text string
							r := regexp.MustCompile(config.Redmine.Url + "([0-9]*)(#note-([0-9]*))?")
							str := r.FindAllStringSubmatch(ev.Text, -1)
							var msg string
							if str != nil {
								log.Printf("str len=%d, %#v\n", len(str), str)
								for i, v := range str {
									log.Printf("str[%d]0=%v\n", i, v[0])
									log.Printf("str[%d]1=%v\n", i, v[1])
									log.Printf("str[%d]2=%v\n", i, v[2])
									text = "redmine " + v[1]
									if v[3] != "" {
										log.Println(v[3])
										no, err := strconv.Atoi(v[3])
										if err != nil {
											log.Println(err)
											break
										}
										m, err := redmineNote(config.Redmine.Url+v[1], no-1)
										if err != nil {
											log.Println(err)
											break
										}
										msg += m
									} else {
										msg += createMessage(text, config.Redmine.Keywords, config.Redmine.Url, redmine)
									}
								}
							}
							postMessage(api, ev, msg)
						} else if ev.User != selfUserId && isWebPage(ev.Text) && isIntraWebPage(ev.Text) {
							var url string
							//if !strings.Contains(ev.Text, config.Redmine.Url) {
							// exludeなら抜ける
							bf := false
							for _, ex := range config.HttpSummary.Exclude {
								if strings.Contains(ev.Text, string(ex.Site)) {
									log.Printf("exclude break: %s, %s\n", ev.Text, string(ex.Site))
									bf = true
									break
								}
							}
							if bf {
								break
							}

							log.Println("not exclude")
							// http Summary
							key := "<(https?://.*." + config.HttpSummary.Intra + "/?.*?)>"
							log.Printf("key=%v\n", key)
							r := regexp.MustCompile(key)
							str := r.FindAllStringSubmatch(ev.Text, -1)
							log.Printf("str=%v\n", str)
							if str != nil {
								log.Printf("str len=%d\n", len(str))
								//var msg string
								for i, v := range str {
									url = v[1]
									log.Printf("str[%d]=%v\n", i, url)
									getImage(api, ev, &url)
								}
								if url != "" {
									postMessage(api, ev, httpSummary(url))
								}
							}
							//}
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

// setLogger Callerを出力
func setLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
