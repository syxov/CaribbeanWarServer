package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"encoding/json"
)

func (self *storage) message(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			self.processError(user, err)
		}
	}()

	var message messagesStructs.Message
	for user.IsInWorld() {
		if err := user.GetConn().ReadJSON(&message); err == nil {
			marshaled, _ := json.Marshal(message)
			switch message.Action {
			case "exitWorld":
				self.remove(user, true)
				return
			case "chat":
				self.chat(&message)
			case "move":
				var moveMessage messagesStructs.MoveIncome
				json.Unmarshal(marshaled, &moveMessage)
				self.move(user, moveMessage)
			case "shoot":
				var shootMessage messagesStructs.ShootIncome
				json.Unmarshal(marshaled, &shootMessage)
				self.shoot(user, shootMessage)
			default:
				user.GetConn().WriteJSON(messagesStructs.ErrorMessage("unrecognized action " + message.Action))
			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user, false)
			} else {
				user.GetConn().WriteJSON(messagesStructs.ErrorMessage(err.Error()))
			}
		}
	}
}

func (self *storage) processError(user *structs.User, err interface{}) {
	var message string
	switch t := err.(type) {
	case error:
		message = t.Error()
	default:
		message = "Something wrong"
	}
	user.GetConn().WriteJSON(messagesStructs.ErrorMessage(message))
	user.GetConn().Close()
	self.remove(user, false)
}
