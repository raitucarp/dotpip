package dotpip

// PubSubMessage represents a message in a pubsub channel.
type PubSubMessage struct {
	Type    string
	Pattern string
	Channel string
	Payload string
}

// PubSubSubscription represents a pubsub subscription.
type PubSubSubscription interface {
	Channel() <-chan PubSubMessage
	Unsubscribe(channels ...string) error
	PUnsubscribe(patterns ...string) error
	SUnsubscribe(shardChannels ...string) error
	Close() error
}
