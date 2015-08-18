package chat

type Message struct {
	Content   string
	Operation string
}

type Origin struct {
	Name    string
	Picture string
}

type OutgoingMessage struct {
	Message
	Origin *Origin
}
