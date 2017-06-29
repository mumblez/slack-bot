package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

const (
	//SHIPPABLE_BOT = "<@devops.bot>" // the bot shippable notifications will be coming from
	//shippableBot = "Shippable"
	shippableBot = "yusuf.tran" // the bot shippable notifications will be coming from
	ciChannel    = "yt-notifications-test"

	//CI_CHANNEL = "ci-notifications" // the channel the bot will message to
)

/*

Fragile parts:
- bot username
- channel
- google spreadsheet
- username changes

*/

func main() {

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			// case *slack.ConnectedEvent:
			// 	fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				fmt.Printf("channel = %+v\n", info.GetChannelByID(ev.Channel).Name)
				// rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, "@yusuf.tran"))

				// act if message came from devops.bot / shippable and if in specific channel
				if info.GetUserByID(ev.User).Name == shippableBot && info.GetChannelByID(ev.Channel).Name == ciChannel {
					process(rtm, ev)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func process(rtm *slack.RTM, msg *slack.MessageEvent) {
	//rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, msg.Channel))
	rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, "@yusuf.tran"))
}
