package structs

import (
	"CaribbeanWarServer/rtree"
	"github.com/gorilla/websocket"
	"math"
	"sync"
)

const (
	left = iota
	right
	none
	angleSpeed         = 10
	velocity   float64 = 1
)

type NearestUser struct {
	ID   *uint           `json:"id"`
	Conn *websocket.Conn `json:""`
	Ship *Ship           `json:"ship"`
	Nick *string         `json:"nick"`
}

type User struct {
	ID                uint   `json:"id"`
	Email             string `json:"email"`
	Nick              string `json:"nick"`
	Cash              uint   `json:"cash"`
	conn              *websocket.Conn
	Location          *Point        `json:"location"`
	Ships             []Ship        `json:"ships"`
	SelectedShip      *Ship         `json:"selectedShip"`
	NearestUsers      []NearestUser `json:"nearestUsers"`
	inWorld           bool
	targetSpeedRatio  float64
	speedRatio        float64
	rotationAngle     float64
	rotationDirection byte
	sync.Mutex
}

func (self *User) Bounds() *rtree.Rect {
	bound, _ := rtree.NewRect(rtree.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	return bound
}

func (self *User) SetMove(moveType string) {
	switch moveType {
	case "upward":
		self.targetSpeedRatio = math.Min(self.targetSpeedRatio+1/3, 1)
	case "backward":
		self.targetSpeedRatio = math.Max(self.targetSpeedRatio-1/3, 0)
	case "left":
		self.rotationDirection = left
	case "right":
		self.rotationDirection = right
	case "none":
		self.rotationDirection = none
	default:
		self.GetConn().WriteJSON(map[string]string{
			"action":  "fuckup",
			"details": "unrecognized command to move" + moveType,
		})
	}
}

func (self *User) UpdatePosition(delta float64) {
	if self.rotationDirection != none {
		if self.rotationDirection == right {
			self.rotationAngle = math.Mod(self.rotationAngle+angleSpeed*delta, math.Pi)
		} else {
			self.rotationAngle = math.Mod(self.rotationAngle-angleSpeed*delta, math.Pi)
		}
	}
	self.speedRatio = lerp(self.speedRatio, self.targetSpeedRatio, delta)
	ship := self.SelectedShip
	if ship != nil {
		self.Location.X += (ship.Speed * self.speedRatio * delta) * math.Cos(self.rotationAngle)
		self.Location.Y += (ship.Speed * self.speedRatio * delta) * math.Sin(self.rotationAngle)
	}
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}

func (self *User) GetConn() *websocket.Conn {
	return self.conn
}

func (self *User) SetConn(conn *websocket.Conn) {
	self.conn = conn
}

func (self *User) IsInWorld() bool {
	return self.inWorld
}

func (self *User) SetIsInWorld(is bool) {
	self.inWorld = is
}
