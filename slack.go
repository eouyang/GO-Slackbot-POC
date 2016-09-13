package main

import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"sync/atomic"
)

type slackSelfResponse struct {
	Id	string	`json:"id"`
	Name	string	`json:"name"`
}

type rtmStartResponse struct {
	Ok	bool	`json:"ok"`
	Url	string	`json:"url"`
	Self	slackSelfResponse `json:"self"`
	Error	string	`json: "error"`
}

func startSlack(token string) (socketUrl string, slackId string, err error) {

	authenticationUrl := "https://slack.com/api/rtm.start?token=" + token;
	resp, err := http.Get(authenticationUrl)
	if err != nil {
		log.Println("Authentication failed", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("Slack API failed")
		err = fmt.Errorf("Slack API failed with status code[" + string(resp.StatusCode) + "]")
		return
	}

	var startSlackResp rtmStartResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&startSlackResp)
	if err != nil {
		return
	}

	if !startSlackResp.Ok {
		log.Println("Slack API failed with error: " + startSlackResp.Error)
		err = fmt.Errorf(startSlackResp.Error)
		return
	}

	return startSlackResp.Url, startSlackResp.Self.Id, err

}

func getSlackConnection(token string) (*websocket.Conn, string, error) {

	socketUrl, id, err := startSlack(token)
	if err != nil {
		return nil, "", err
	}
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		return nil, "", err
	}
	return conn, id, err
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func postMessageToSlack(conn *websocket.Conn, m Message, counter *uint64) (err error) {
	m.Id = atomic.AddUint64(counter, 1)
	return conn.WriteJSON(m)
}

func readMessageFromSlack(conn *websocket.Conn) (m Message, err error) {
	err = conn.ReadJSON(&m)
	return m, err

}

