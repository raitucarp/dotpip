package dotpip

type PubSubMessage struct {
	Type    string
	Pattern string
	Channel string
	Payload string
}

type PubSubSubscription interface {
	Channel() <-chan PubSubMessage
	Unsubscribe(channels ...string) error
	PUnsubscribe(patterns ...string) error
	SUnsubscribe(shardChannels ...string) error
	Close() error
}
