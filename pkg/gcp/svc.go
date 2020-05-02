package svc

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type GCPServiceProvider struct {
	GCPClient *pubsub.Client
}

func NewGCPServiceProvider(id, region string) *GCPServiceProvider {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, id)
	if err != nil {
		log.Println("Error in client creation", err)
	}
	return &GCPServiceProvider{
		GCPClient: client,
	}
}

func (p *GCPServiceProvider) PublishEvent(subject, region string, eventData interface{}) {
	ctx := context.Background()
	t := p.GCPClient.Topic(subject)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(eventData.(string)),
	})
	res, err := result.Get(ctx)
	if err != nil {
		log.Println("Publish error", err)
	} else {
		log.Println("message published", res)
	}
	t.Stop()
}

func (p *GCPServiceProvider) SubscribeEvent(subject, region string, onEventHandler func(interface{})) {
	ctx := context.Background()

	t := p.GCPClient.Topic(subject)
	s := p.GCPClient.Subscription(subject + "_Q")

	ok, err := s.Exists(ctx)
	if err != nil {
		log.Println("Error in checking subscription", err)
	}

	if !ok {
		s, err = p.GCPClient.CreateSubscription(context.Background(), subject+"_Q", pubsub.SubscriptionConfig{Topic: t})
		if err != nil {
			log.Println("Error in creating subscription", err)
			return
		}
	}

	s.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
		onEventHandler(msg.Data)
		msg.Ack()
	})
}
