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
			user.GetConn().WriteJSON(errors.New(message))
			user.GetConn().Close()
			self.remove(user, false)
		}
	}()

	var json map[string]interface{}
	for user.IsInWorld() {
		if err := user.GetConn().ReadJSON(&json); err == nil {
			details := json["details"].(map[string]interface{})
			action := json["action"].(string)
			switch action {
			case "exitWorld":
				self.remove(user, true)
				return
			case "chat":
				self.chat(&json)
			case "move":
				self.move(user, details)
			case "shoot":
				self.shoot(user, details)
			default:
				user.GetConn().WriteJSON(map[string]string{
					"action":  "fuckup",
					"details": "unrecognizedAction " + action,
				})

			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user, false)
				return
			} else {
				user.GetConn().WriteJSON(map[string]string{
					"action":  "fuckup",
					"details": err.Error(),
				})
			}
		}
	}
}
