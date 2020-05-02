package main

import (
	"eventus/pkg/gcp"
	"eventus/pkg/provider"
	"log"
)

func OnEventHandler(eventData interface{}) {
	switch eventData.(type) {
	case []uint8:
		log.Println("Message:", string(eventData.([]uint8)))
	case string:
		log.Println("Message:", eventData.(string))
	default:
		log.Println("Message data need type handling")
	}
}

func main() {
	log.Println("Event wrapper library in GO")
	providerGCP := provider.CustomProvider{
		PubSubModule: svc.NewGCPServiceProvider("origin-1205", "us-west-2"),
	}

	providerGCP.PublishEvent("genesis", "us-west-2", "Hello world !!!")
	providerGCP.SubscribeEvent("genesis", "us-west-2", OnEventHandler)
}
