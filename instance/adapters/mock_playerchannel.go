package adapters

import "github.com/bensivo/salad-bowl/instance"

// PlayerChannel implementation used for testing
// Simply saves a record of each sent message, and exposes a function for mimicking a message sent from the player to the instance.
type MockPlayerChannel struct {
	Sent      []interface{}
	Callbacks []instance.MessageCallback
}

var _ instance.PlayerChannel = (*MockPlayerChannel)(nil)

func NewMockPlayerChannel() *MockPlayerChannel {
	return &MockPlayerChannel{
		Sent:      []interface{}{},
		Callbacks: []instance.MessageCallback{},
	}
}

// Send a message to the player
func (tpc *MockPlayerChannel) Send(message interface{}) error {
	tpc.Sent = append(tpc.Sent, message)
	return nil
}

// Simulate a message sent from the player to the instance
func (tpc *MockPlayerChannel) MockReceive(message interface{}) {
	for _, cb := range tpc.Callbacks {
		cb(message)
	}
}

func (tpc *MockPlayerChannel) OnMessage(cb instance.MessageCallback) {
	tpc.Callbacks = append(tpc.Callbacks, cb)
}
