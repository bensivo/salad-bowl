package game

type MessageCallback func(message interface{}) error

// PlayerChannel is any channel which can be used to communicate with a player device.
// This interface allows different players to connect through different mechanisms.
//
// PlayerChannel implementations must be bidirectional
type PlayerChannel interface {
	Send(message interface{}) error

	OnMessage(cb MessageCallback)

	// TODO add a standard interface for communicating health (online/offline, bandwidth, up to date?)
}
