package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

const (
	//SHIPPABLE_BOT = "<@devops.bot>" // the bot shippable notifications will be coming from
	shippableBot = "Shippable"
	//shippableBot = "yusuf.tran" // the bot shippable notifications will be coming from
	ciChannel = "yt-notifications-test"

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
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				// todo : add logic to grab user mappings and search for user their first
				users, err := rtm.GetUsers()
				if err != nil {
					fmt.Printf("Error getting users: %v\n", err)
				}

				for _, user := range users {
					// amend logic to lookup github and slack user mapping!
					if user.Name == "yusuf.tran" && ev.Username == shippableBot && info.GetChannelByID(ev.Channel).Name == ciChannel {
						_, _, chat, err := rtm.OpenIMChannel(user.Name)
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

func process(rtm *slack.RTM, msg *slack.MessageEvent) {
	//rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, msg.Channel))
	// rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, "yusuf.tran"))
	// rtm.SendMessage(rtm.NewOutgoingMessage(msg.Text, "@yusuf.tran"))
	// params := &slack.PostMessageParameters{}
	_, _, chat, err := rtm.OpenIMChannel("yusuf.tran")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	rtm.SendMessage(rtm.NewOutgoingMessage("123...", chat))

}
