package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data map[string]interface{}) {
	moveType := data["type"].(string)
	self.Lock()
	user.SetMove(moveType)
	self.Unlock()
	user.Lock()
	for _, value := range user.NearestUsers {
		value.Conn.WriteJSON(map[string]interface{}{
			"action": "move",
			"details": map[string]interface{}{
				"id":   user.ID,
				"type": moveType,
			},
		})
	}
	user.Unlock()
}

func (self *storage) movement(user *structs.User) {
	defer func() {
		if recover() != nil {
			self.Unlock()
		}
	}()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for user.InWorld {
		time := <-ticker.C
		self.Lock()
		self.ocean.Delete(user)
		user.UpdatePosition(float64(time.Nanosecond() / 1000))
		self.ocean.Insert(user)
		self.Unlock()
	}
}
