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
	chatCn := make(chan *messagesStructs.Message, 5)
	moveCh := make(chan *messagesStructs.MoveIncome, 5)
	shootCh := make(chan *messagesStructs.ShootIncome, 5)
	go self.chat(chatCn)
	go self.move(user, moveCh)
	go self.shoot(user, shootCh)
	for user.IsInWorld() {
		if err := user.GetConn().ReadJSON(&message); err == nil {
			marshaled, _ := json.Marshal(message)
			switch message.Action {
			case "exitWorld":
				go self.remove(user, true)
				break
			case "chat":
				chatCn <- &message
			case "move":
				var moveMessage messagesStructs.MoveIncome
				json.Unmarshal(marshaled, &moveMessage)
				moveCh <- &moveMessage
			case "shoot":
				var shootMessage messagesStructs.ShootIncome
				json.Unmarshal(marshaled, &shootMessage)
				shootCh <- &shootMessage
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
	close(chatCn)
	close(moveCh)
	close(shootCh)
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
