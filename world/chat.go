package world

func (self *WorldStruct) chat(id uint, data map[string]interface{}) {
	sender := data["details"].(map[string]string)["sender"]
	userMessage := data["details"].(map[string]string)["message"]
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  sender,
			"message": userMessage,
		},
	}
	for k, v := range self.world {
		if k != id {
			v.conn.WriteJSON(message)
		}
	}
}
