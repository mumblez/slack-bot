package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

const (
	userMappingsFile = "./userMappings.json"
	shippableBot     = "Shippable"
	ciChannel        = "yt-notifications-test"
	//CI_CHANNEL = "ci-notifications" // the channel the bot will message to
)

type githubToSlack struct {
	Github string `json:"github"`
	Slack  string `json:"slack"`
}

type users struct {
	Collection []githubToSlack `json:"github_to_slack"`
}

func findUser(user string, userMap *users) (slackUser string, ok bool) {
	for _, u := range userMap.Collection {
		if u.Github == user {
			return u.Slack, true
		}
	}
	return "", false
}

func main() {

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	file, err := ioutil.ReadFile(userMappingsFile)
	if err != nil {
		fmt.Printf("Error reading userMappings.json file = %+v\n", err)
		os.Exit(1)
	}

	userMap := &users{Collection: make([]githubToSlack, 0)}
	if err := json.Unmarshal(file, &userMap); err != nil {
		fmt.Printf("err = %+v\n", err)
	}

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				// todo : extract from text, add logic to grab user mappings and search for user their first

				// extract github username from message (2nd line)
				slackMsg := strings.Split(ev.Attachments[0].Text, "\n")
				msg := strings.Fields(slackMsg[1])
				githubUser := msg[len(msg)-1]

				// lookup and check github to slack username mapping
				slackUser, ok := findUser(githubUser, userMap)
				if !ok {
					fmt.Printf("githubUser not found = %+v\n", githubUser)
					continue
				}

				users, err := rtm.GetUsers()
				if err != nil {
					fmt.Printf("Error getting users: %v\n", err)
				}

				for _, user := range users {
					// amend logic to lookup github and slack user mapping!
					if user.Name == slackUser && ev.Username == shippableBot && info.GetChannelByID(ev.Channel).Name == ciChannel {
						_, _, chat, err := rtm.OpenIMChannel(user.ID)
						if err != nil {
							fmt.Printf("Error opening channel to user: %v\n", err)
						}
						// process(rtm, ev)
						params := &slack.PostMessageParameters{
							AsUser:      true,
							Attachments: ev.Attachments,
						}
						_, _, err = rtm.PostMessage(chat, ev.Text, *params)
						if err != nil {
							fmt.Printf("Error posting message to user: %v\n", err)
						}
					}
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}
