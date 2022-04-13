package main

import (
	"flag"
	"log"
	tgClient "remindbot/clients/telegram"
	event_consumer "remindbot/consumer/event-consumer"
	"remindbot/events/telegram"
	"remindbot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)
	log.Print("service has been started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service has been stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"telegram access token",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}

	return *token
}
