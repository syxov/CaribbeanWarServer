package messagesStructs

type Message struct {
	Action  string      `json:"action"`
	Details interface{} `json:"details"`
}

type errorMessage struct {
	Message string `json:"message"`
}

func ErrorMessage(message string) Message {
	return Message{"error", errorMessage{message}}
}
