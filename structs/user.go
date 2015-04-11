package structs

import (
	"CaribbeanWarServer/commonStructs"
	"CaribbeanWarServer/intmath"
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

type User struct {
	ID                uint                       `json:"id"`
	Email             string                     `json:"email"`
	Nick              string                     `json:"nick"`
	Cash              uint                       `json:"cash"`
	Location          point.Point2D              `json:"location"`
	Ships             []commonStructs.Ship       `json:"ships"`
	SelectedShip      *commonStructs.Ship        `json:"selectedShip"`
	NearestUsers      commonStructs.NearestUsers `json:"nearestUsers"`
	RotationAngle     float64                    `json:"alpha"`
	conn              *commonStructs.Connection
	inWorld, killed   atomic.Value
	sailsMode         int32
	speedRatio        float64
	rotationDirection int32
	sync.Mutex
}

func (self *User) Bounds(radius ...float64) *rtree.Rect {
	var value float64 = 5
	if len(radius) != 0 {
		value = radius[0]
	}
	return rtree.NewRect(rtree.Point{self.Location.X - value/2, self.Location.Y - value/2}, []float64{value, value}, self.RotationAngle)
}

func (self *User) SetMove(moveType string) {
	switch moveType {
	case "upward":
		atomic.StoreInt32(&self.sailsMode, intmath.Min32(self.sailsMode+1, 3))
	case "backward":
		atomic.StoreInt32(&self.sailsMode, intmath.Max32(self.sailsMode-1, 0))
	case "left":
		atomic.StoreInt32(&self.rotationDirection, left)
	case "right":
		atomic.StoreInt32(&self.rotationDirection, right)
	case "none":
		atomic.StoreInt32(&self.rotationDirection, none)
	default:
		self.GetConn().WriteJSON(messagesStructs.ErrorMessage("ERRORS_UNKNOWN_ACTION"))
	}
}

func (self *User) UpdatePosition() {
	self.Lock()
	ship := self.SelectedShip
	if ship != nil {
		self.Location.X += self.speedRatio * math.Cos(self.RotationAngle)
		self.Location.Y += self.speedRatio * math.Sin(-self.RotationAngle)
		rotationDirection := float64(atomic.LoadInt32(&self.rotationDirection))
		sailsMode := float64(atomic.LoadInt32(&self.sailsMode))
		self.RotationAngle = math.Mod(self.RotationAngle+(rotationDirection*angleSpeed*self.speedRatio)/(sailsMode+1.0), 2*math.Pi)
	}
	self.Unlock()
}

func (self *User) SendForAll(message interface{}) {
	self.Lock()
	self.conn.WriteJSON(message)
	for _, p := range self.NearestUsers {
		p.Conn.WriteJSON(message)
	}
	self.Unlock()
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

func (self *User) IsKilled() bool {
	return self.killed.Load().(bool)
}

func (self *User) SetIsKilled(is bool) {
	self.killed.Store(is)
}

func (self *User) UpdateSpeed(delta float64) float64 {
	if ship := self.SelectedShip; ship != nil {
		self.speedRatio = lerp(self.speedRatio, float64(atomic.LoadInt32(&self.sailsMode))*ship.Speed*delta/4.0, velocity)
	}
	return self.speedRatio
}

func (self *User) DoKill() {
	self.SetIsKilled(true)
	self.speedRatio = 0
	self.sailsMode = 0
	self.rotationDirection = 0
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}
