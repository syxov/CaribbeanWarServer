package world

import (
	"CaribbeanWarServer/structs"
	"errors"
)

func (self *storage) message(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			var message string
			switch t := err.(type) {
			case error:
				message = t.Error()
			default:
				message = "Something wrong"
			}
			user.Conn.WriteJSON(errors.New(message))
			user.Conn.Close()
			self.remove(user)
		}
	}()

	var message interface{}
	for user.InWorld {
		if err := user.Conn.ReadJSON(&message); err == nil {
			json := message.(map[string]interface{})
			details := json["details"].(map[string]interface{})
			switch json["action"] {
			case "exitWorld":
				self.remove(user)
				return
			case "chat":
				self.chat(&message)
			case "move":
				self.move(user, details)
			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user)
				return
			} else {
				user.Conn.WriteJSON(map[string]string{
					"action":  "fuckup",
					"details": err.Error(),
				})
			}
		}
	}
}
