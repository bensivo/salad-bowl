package adapters

import "github.com/bensivo/salad-bowl/hub"

// PlayerChannel implementation used for testing
// Simply saves a record of each sent message, and exposes a function for mimicking a message sent from the player to the hub.
type MockPlayerChannel struct {
	Sent      []interface{}
	Callbacks []hub.MessageCallback
}

var _ hub.PlayerChannel = (*MockPlayerChannel)(nil)

func NewMockPlayerChannel() *MockPlayerChannel {
	return &MockPlayerChannel{
		Sent:      []interface{}{},
		Callbacks: []hub.MessageCallback{},
	}
}

// Send a message to the player
func (tpc *MockPlayerChannel) Send(message interface{}) error {
	tpc.Sent = append(tpc.Sent, message)
	return nil
}

// Simulate a message sent from the player to the hub
func (tpc *MockPlayerChannel) MockReceive(message interface{}) {
	for _, cb := range tpc.Callbacks {
		cb(message)
	}
}

func (tpc *MockPlayerChannel) OnMessage(cb hub.MessageCallback) {
	tpc.Callbacks = append(tpc.Callbacks, cb)
}
