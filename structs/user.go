package structs

import (
	"CaribbeanWarServer/commonStructs"
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/rtree"
	"math"
	"sync"
	"sync/atomic"
)

const (
	left       = -1
	right      = 1
	none       = 0
	angleSpeed = 0.075
	velocity   = 0.01
)

type sailsModeType int8
type rotationType int8

type User struct {
	ID                uint                        `json:"id"`
	Email             string                      `json:"email"`
	Nick              string                      `json:"nick"`
	Cash              uint                        `json:"cash"`
	Location          point.Point2D               `json:"location"`
	Ships             []commonStructs.Ship        `json:"ships"`
	SelectedShip      *commonStructs.Ship         `json:"selectedShip"`
	NearestUsers      []commonStructs.NearestUser `json:"nearestUsers"`
	RotationAngle     float64                     `json:"alpha"`
	conn              *commonStructs.Connection
	inWorld           atomic.Value
	sailsMode         sailsModeType
	speedRatio        float64
	rotationDirection rotationType
	sync.Mutex
}

func (self *User) Bounds(radius ...float64) *rtree.Rect {
	var value float64 = 5
	if len(radius) != 0 {
		value = radius[0]
	}
	return rtree.NewRect(rtree.Point{self.Location.X - value/2, self.Location.Y - value/2}, []float64{value, value})
}

func (self *User) SetMove(moveType string) {
	self.Lock()
	switch moveType {
	case "upward":
		self.sailsMode = sailsModeType(math.Min(float64(self.sailsMode+1), 3))
	case "backward":
		self.sailsMode = sailsModeType(math.Max(float64(self.sailsMode-1), 0))
	case "left":
		self.rotationDirection = left
	case "right":
		self.rotationDirection = right
	case "none":
		self.rotationDirection = none
	default:
		self.GetConn().WriteJSON(messagesStructs.ErrorMessage("ERRORS_UNKNOWN_ACTION"))
	}
	self.Unlock()
}

func (self *User) UpdatePosition(delta float64) {
	self.Lock()
	ship := self.SelectedShip
	if self.SelectedShip != nil {
		self.speedRatio = lerp(self.speedRatio, float64(self.sailsMode)*ship.Speed*delta/4.0, velocity)
		self.Location.X += self.speedRatio * math.Cos(self.RotationAngle)
		self.Location.Y += self.speedRatio * math.Sin(-self.RotationAngle)
		self.RotationAngle = math.Mod(self.RotationAngle+(float64(self.rotationDirection)*angleSpeed*self.speedRatio)/(float64(self.sailsMode)+1.0), 2*math.Pi)
		self.GetConn().WriteJSON(messagesStructs.Message{"position", map[string]interface{}{
			"x":     self.Location.X,
			"y":     self.Location.Y,
			"alpha": self.RotationAngle,
		}})
	}
	self.Unlock()
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}

func (self *User) GetConn() *commonStructs.Connection {
	return self.conn
}

func (self *User) SetConn(conn *commonStructs.Connection) {
	self.conn = conn
}

func (self *User) IsInWorld() bool {
	return self.inWorld.Load().(bool)
}

func (self *User) SetIsInWorld(is bool) {
	self.inWorld.Store(is)
}
