package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) shoot(user *structs.User, details map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			user.GetConn().WriteJSON(map[string]interface{}{
				"action": "error",
				"details": map[string]interface{}{
					"from":    "shoot",
					"message": err,
				},
			})
		}
	}()
	position := details["location"].(map[string]interface{})
	angle := details["angle"].(float64)
	direction := details["direction"].(float64)
	user.Lock()
	core := structs.NewCore(&structs.Point3D{position["x"].(float64), position["y"].(float64), position["z"].(float64)}, angle, direction, user.ID)
	message := map[string]interface{}{
		"action": "shoot",
		"details": map[string]interface{}{
			"id":       user.ID,
			"alpha":    user.RotationAngle,
			"angle":    angle,
			"location": position,
		},
	}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(message)
	}
	user.Unlock()
	user.GetConn().WriteJSON(message)
	go self.updateCore(core, user)
}

func (self *storage) updateCore(core *structs.Core, user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			user.GetConn().WriteJSON(map[string]interface{}{
				"action": "error",
				"details": map[string]interface{}{
					"from":    "shoot",
					"message": err.(error).Error(),
				},
			})
		}
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for !core.UnderWater() {
		now := time.Now().UnixNano()
		<-timer.C
		core.UpdatePosition(float64(time.Now().UnixNano()-now) / float64(time.Second))
		self.Lock()
		spatials := self.ocean.SearchIntersect(core.GetBounds())
		if len(spatials) > 0 {
			looser := spatials[0].(*structs.User)
			looser.Lock()
			for _, neigbour := range looser.NearestUsers {
				neigbour.Conn.WriteJSON(map[string]interface{}{
					"action": "hit",
					"details": map[string]interface{}{
						"id":       looser.ID,
						"location": core.CurrentPosition,
						"damage":   87,
					},
				})
			}
			looser.Unlock()
			self.Unlock()
			return
		}
		self.Unlock()
	}
}
