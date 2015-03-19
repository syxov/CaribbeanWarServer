package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) shoot(user *structs.User, details map[string]interface{}) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	position := details["location"].(map[string]interface{})
	angle := details["angle"].(float64)
	direction := details["direction"].(float64)
	user.Lock()
	message := structs.Message{"shoot", map[string]interface{}{
		"id":       user.ID,
		"alpha":    user.RotationAngle,
		"angle":    angle,
		"location": position,
	}}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(message)
	}
	user.Unlock()
	user.GetConn().WriteJSON(message)
	core := structs.NewCore(&structs.Point3D{position["x"].(float64), position["y"].(float64), position["z"].(float64)}, angle, direction, user.ID)
	go self.updateCore(core, user)
}

func (self *storage) updateCore(core *structs.Core, user *structs.User) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for !core.UnderWater() {
		now := time.Now().UnixNano()
		<-timer.C
		core.UpdatePosition(float64(time.Now().UnixNano()-now) / float64(time.Second))
		spatials := self.ocean.SearchIntersectWithLimit(1, core.GetBounds())
		if len(spatials) == 1 {
			looser := spatials[0].(*structs.User)
			looser.Lock()
			message := structs.Message{"hit", map[string]interface{}{
				"id":       looser.ID,
				"location": core.CurrentPosition,
				"damage":   87,
			}}
			for _, neigbour := range looser.NearestUsers {
				neigbour.Conn.WriteJSON(message)
			}
			looser.Unlock()
			return
		}
	}
}

func sendErrorMessage(user *structs.User, err interface{}) {
	if err != nil {
		var message = "Smth wrong: shoot"
		switch smth := err.(type) {
		case error:
			message = smth.Error()
		case string:
			message = smth
		}
		user.GetConn().WriteJSON(structs.ErrorMessage(message))
	}
}
