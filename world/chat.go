package world

func (self *WorldStruct) chat(id uint, data map[string]interface{}) {
	convertedData := data["details"].(map[string]interface{})
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  convertedData["sender"],
			"message": convertedData["message"],
		},
	}
	for k, v := range self.world {
		if k != id {
			v.conn.WriteJSON(message)
		}
	}
}
