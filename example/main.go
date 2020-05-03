package main

import (
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

	//region can be passed as empty string for GCP
	providerGCP := provider.NewCustomProvider("gcp", "Account-ID/Project-ID", "region")

	providerGCP.PublishEvent("topic-name", "region", "Hello world !!!")
	providerGCP.SubscribeEvent("topic-name", "region", OnEventHandler)
}
