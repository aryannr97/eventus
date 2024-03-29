package provider

import (
	"sync"

	svc "github.com/aryannr97/eventus/pkg/aws"
	svc2 "github.com/aryannr97/eventus/pkg/gcp"
	"github.com/aryannr97/eventus/pkg/pubsub"
)

var wg sync.WaitGroup

type CustomProvider struct {
	PubSubModule pubsub.EventStore
}

func NewCustomProvider(providerKey, id, region string) *CustomProvider {
	switch providerKey {
	case "aws":
		return &CustomProvider{
			PubSubModule: svc.NewAWSServiceProvider(id, region),
		}
	case "gcp":
		return &CustomProvider{
			PubSubModule: svc2.NewGCPServiceProvider(id, region),
		}
	}
	return nil
}

func (p *CustomProvider) PublishEvent(subject, region string, eventData interface{}) {
	p.PubSubModule.PublishEvent(subject, region, eventData)
}

func (p *CustomProvider) SubscribeEvent(subject, region string, onEventHandler func(interface{})) {
	wg.Add(1)
	go p.PubSubModule.SubscribeEvent(subject, region, onEventHandler)
	wg.Wait()
}
