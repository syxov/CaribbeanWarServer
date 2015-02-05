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

	var json map[string]interface{}
	for user.InWorld {
		if err := user.Conn.ReadJSON(&json); err == nil {
			details := json["details"].(map[string]interface{})
			switch json["action"] {
			case "exitWorld":
				self.remove(user)
				return
			case "chat":
				self.chat(&json)
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
