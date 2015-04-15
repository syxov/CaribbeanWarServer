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
	user.SendForAll(messagesStructs.RespawnOutcome{
		Message: messagesStructs.Message{
			Action: "respawn",
		},
		Details: messagesStructs.RespawnOutcomeDetails{
			ID:       user.ID,
			Location: user.Location,
			Rotation: user.RotationAngle,
		},
	})
}
