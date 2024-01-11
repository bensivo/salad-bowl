package hub

type Message struct {
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}
