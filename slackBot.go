package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"strconv"
)

/* Slack Token */
const authToken = "YOUR SLACK BOT TOKEN"

var phraseMap map[string]int = make(map[string]int);

type SlackBot struct {
	Id	string
	MessageCount	uint64
	SlackConnection	*websocket.Conn
}

func (bot *SlackBot) run() {
	for {
		message, err := readMessageFromSlack(bot.SlackConnection)
		if err != nil {
			log.Fatal(err)
		}

		if message.Type == "message" && strings.HasPrefix(message.Text, "<@"+bot.Id+">") {
			tokens := strings.Fields(message.Text)
			if len(tokens) == 3 && strings.ToUpper(tokens[1]) == "TRACK" {
				if count, ok := phraseMap[strings.ToUpper(tokens[2])]; ok == true {
					message.Text = "```Already tracking [" + strings.ToUpper(tokens[2]) + "]. Phrase has occurred [" + strconv.Itoa(count) + "] times in recorded conversations.```"
				} else {
					phraseMap[strings.ToUpper(tokens[2])] = 0
					message.Text = "```Now tracking [" + strings.ToUpper(tokens[2]) + "]```"
				}
			} else if len(tokens) == 4 && strings.ToUpper(tokens[1]) == "STOP" && strings.ToUpper(tokens[2]) == "TRACKING" {
				if _, ok := phraseMap[strings.ToUpper(tokens[3])]; ok == true {
					delete(phraseMap, strings.ToUpper(tokens[3]))
					message.Text = "```[" + strings.ToUpper(tokens[3]) + "] has been removed from list of tracked phrases```"
				} else {
					message.Text = "``` [" + strings.ToUpper(tokens[3]) + "] is not currently being tracked```"
				}
			} else {
				message.Text = "```Usage:\n1.track [phrase]\n2.stop tracking [phrase]```"
			}
			postMessageToSlack(bot.SlackConnection, message, &bot.MessageCount)
		} else {
			tokens := strings.Fields(message.Text)
			for _, phrase := range tokens {
				if count, ok := phraseMap[strings.ToUpper(phrase)]; ok == true {
					phraseMap[strings.ToUpper(phrase)] = count + 1
				}
			}
		}
	}
}

func (bot *SlackBot) init() {
	conn, id, err := getSlackConnection(authToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Id = id
	bot.SlackConnection = conn

	fmt.Println("Slack Bot ready to serve...")
}