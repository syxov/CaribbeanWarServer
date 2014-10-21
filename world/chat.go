package world

func (self *storage) chat(data map[string]interface{}) {
	message := map[string]interface{}{
		"Action": "chat",
		"Details": map[string]interface{}{
			"Sender":  data["Sender"],
			"Message": data["Message"],
		},
	}
	self.Lock()
	defer self.Unlock()
	for _, v := range self.userList {
		v.Conn.WriteJSON(message)
	}
}
