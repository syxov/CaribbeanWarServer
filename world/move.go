package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data map[string]interface{}) {
	moveType := data["type"].(string)
	sendData := map[string]interface{}{
		"action": "move",
		"details": map[string]interface{}{
			"id":   user.ID,
			"type": moveType,
		},
	}
	user.Lock()
	user.SetMove(moveType)
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.Unlock()
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for user.InWorld {
		tick := <-ticker.C
		self.Lock()
		self.ocean.Delete(user)
		user.UpdatePosition(float64(tick.Nanosecond() / int(time.Second)))
		self.ocean.Insert(user)
		self.Unlock()
	}
}
