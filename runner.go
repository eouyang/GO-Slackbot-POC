package main

import "fmt"

func main() {
	fmt.Println("Starting slack bot")
	var bot SlackBot
	bot.init()
	bot.run()
}
