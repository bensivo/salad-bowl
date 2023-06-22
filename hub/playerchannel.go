package hub

type MessageCallback func(message interface{}) error

// PlayerChannel is any interface or protocol which can be used to communicate with a player device.
type PlayerChannel interface {
	Send(message interface{}) error

	OnMessage(cb MessageCallback)
}
