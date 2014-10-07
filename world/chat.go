package world

func (self *WorldStruct) chat(data map[string]interface{}) {
	sender := data["details"].(map[string]uint)["sender"]
	userMessage := data["details"].(map[string]string)["message"]
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  sender,
			"message": userMessage,
		},
	}
	for k, v := range self.world {
		if k != sender {
			v.conn.WriteJSON(message)
		}
	}
}
