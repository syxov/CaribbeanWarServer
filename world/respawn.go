package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/structs"
)

func (self *storage) respawn(user *structs.User, ch chan bool) {
	for {
		if _, ok := <-ch; ok {
			self.doRespawn(user)
		} else {
			return
		}
	}
}

func (self *storage) doRespawn(user *structs.User) {
	user.Location = point.Point2D{0, 0}
	user.RotationAngle = 0.3
	user.SetIsKilled(false)
	user.GetConn().WriteJSON(messagesStructs.RespawnOutcome{
		Details: messagesStructs.RespawnOutcomeDetails{
			Location: user.Location,
			Rotation: user.RotationAngle,
		},
		Message: messagesStructs.Message{
			Action: "respawn",
		},
	})
}
