package pubsub

type EventStore interface {
	PublishEvent(string, string, interface{})
	SubscribeEvent(string, string, func(interface{}))
}
