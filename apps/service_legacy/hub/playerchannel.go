package hub

type MessageCallback func(message Message)

type DisconnectCallback func()

// PlayerChannel is any interface or protocol which can be used to communicate with a player device.
//
//go:generate mockery --name PlayerChannel
type PlayerChannel interface {
	Send(message interface{}) error

	// OnMessages adds a callback for messages received from this channel
	OnMessage(cb MessageCallback)

	// OnDisconnect adds a callback for disconnect events on this channel
	OnDisconnect(cb DisconnectCallback)

	Close() error
}
