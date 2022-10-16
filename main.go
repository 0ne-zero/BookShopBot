package main

import (
	"log"
)

func main() {
	// Config logger
	f, err := openLogFile()
	if err != nil {
		log.Fatalf("Error occurred during open log file - %s\n", err.Error())
	}
	log.SetOutput(f)
	log.SetFlags(log.Llongfile)

	log.Println("Starting ...")

	// Create bot
	// Config how to update messages
	bot, updates, err := configBot()
	if err != nil {
		log.Fatalf("Error occurred during config bot - %s\n", err.Error())
	}

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	for update := range updates {
		if isCommand(update.Message.Text) {
			switch update.Message.Text {
			case "/start":

			}
		}
	}
}
