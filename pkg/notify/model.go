package notify

type Message struct {
	Heading string            `json:"heading"`
	Message string            `json:"message"`
	Payload map[string]string `json:"payload"`
}
