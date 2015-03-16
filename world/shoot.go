package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) shoot(user *structs.User, details map[string]interface{}) {
	position := details["position"].(map[string]float64)
	angle := details["angle"].(float64)
	core := structs.NewCore(&structs.Point3D{position["x"], position["y"], position["z"]}, angle, user.RotationAngle, user.ID)
	message := map[string]interface{}{
		"action": "shoot",
		"details": map[string]interface{}{
			"id":       user.ID,
			"alpha":    user.RotationAngle,
			"angle":    angle,
			"location": position,
		},
	}
	user.Lock()
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(message)
	}
	user.GetConn().WriteJSON(message)
	user.Unlock()
	go self.updateCore(core)
}

func (self *storage) updateCore(core *structs.Core) {
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for core.UnderWater() {
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
