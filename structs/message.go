package structs

type Message struct {
	Action  string                 `json:"action"`
	Details map[string]interface{} `json:"details"`
}

func ErrorMessage(data string) Message {
	return Message{"error", map[string]interface{}{
		"message": data,
	}}
}
