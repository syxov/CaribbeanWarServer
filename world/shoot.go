package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"github.com/syxov/intmath"
	"time"
)

func (self *storage) shoot(user *structs.User, ch chan *messagesStructs.ShootIncome) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	for {
		if incomeMessage, ok := <-ch; ok {
			details := incomeMessage.Details
			user.Lock()
			message := messagesStructs.ShootOutcome{
				Action: "shoot",
				Details: messagesStructs.ShootOutcomeDetails{
					ID:        user.ID,
					Alpha:     user.RotationAngle,
					Angle:     details.Angle,
					Location:  details.Location,
					Direction: details.Direction,
				},
			}
			core := structs.NewCore(details.Location, user.RotationAngle, details.Angle, details.Direction, user.ID)
			user.Unlock()
			go self.updateCore(core, user)
			go user.SendForAll(message)
		} else {
			return
		}
	}
}

func (self *storage) updateCore(core *structs.Core, user *structs.User) {
	timer := time.NewTicker(10 * time.Millisecond)
	defer func() {
		timer.Stop()
		sendErrorMessage(user, recover())
	}()
	for core.OverWater() {
		now := time.Now().UnixNano()
		<-timer.C
		core.UpdatePosition(float64(time.Now().UnixNano()-now) / float64(time.Second))
		spatials := self.ocean.SearchIntersectWithLimit(1, core.GetBounds(), func(spat *rtree.Spatial) bool {
			spatialUser := (*spat).(*structs.User)
			return spatialUser.ID != user.ID && !spatialUser.IsKilled()
		})
		if len(spatials) != 0 {
			looser := spatials[0].(*structs.User)
			message := messagesStructs.Hit{
				Action: "hit",
				Details: messagesStructs.HitDetails{
					ID:       looser.ID,
					Location: core.CurrentPosition,
					Damage:   87,
				},
			}
			looser.SendForAll(message)
			looser.Lock()
			looser.SelectedShip.CurrentHP = intmath.Max(looser.SelectedShip.CurrentHP-87, 0)
			if looser.SelectedShip.CurrentHP == 0 {
				looser.DoKill()
				message := messagesStructs.Dead{
					Action: "death",
					Details: messagesStructs.DeadDetails{
						ID:       looser.ID,
						Location: looser.Location,
						Rotation: looser.RotationAngle,
					},
				}
				looser.Unlock()
				looser.SendForAll(message)
			} else {
				looser.Unlock()
			}
			return
		}
	}
	missMessage := messagesStructs.Miss{
		Message: messagesStructs.Message{
			Action: "miss",
		},
		Details: messagesStructs.MissDetails{
			Location: core.CurrentPosition,
		},
	}
	user.SendForAll(missMessage)
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
		user.GetConn().WriteJSON(messagesStructs.ErrorMessage(message))
	}
}
