package world

import (
	"CaribbeanWarServer/structs"
	"errors"
)

func (self *storage) message(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			self.processError(user, err)
		}
	}()

	var message structs.Message
	for user.IsInWorld() {
		if err := user.GetConn().ReadJSON(&message); err == nil {
			switch message.Action {
			case "exitWorld":
				self.remove(user, true)
				return
			case "chat":
				self.chat(&message)
			case "move":
				self.move(user, message.Details)
			case "shoot":
				self.shoot(user, message.Details)
			default:
				user.GetConn().WriteJSON(structs.ErrorMessage("unrecognized action " + message.Action))
			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user, false)
			} else {
				user.GetConn().WriteJSON(structs.ErrorMessage(err.Error()))
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
	user.GetConn().WriteJSON(errors.New(message))
	user.GetConn().Close()
	self.remove(user, false)
}
