package structs

import (
	"CaribbeanWarServer/rtree"
	"math"
	"sync"
	"github.com/gorilla/websocket"
)

const (
	left = iota
	right
	none
	angleSpeed         = 10
	velocity   float64 = 1
)

type NearestUser struct {
	ID   uint
	Conn *websocket.Conn
	Ship *Ship
	Nick string
}

type User struct {
	ID                uint
	Email             string
	Nick              string
	Cash              uint
	Conn              *websocket.Conn
	Location          *Point
	Ships             []Ship
	SelectedShip      *Ship
	NearestUsers      []NearestUser
	InWorld           bool
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
		self.Conn.WriteJSON(map[string]string{
			"action":  "fuckup",
			"details": "unrecognized command to move" + moveType,
		})
	}
}

func (self *User) UpdatePosition(delta float64) {
	self.Lock()
	if self.rotationDirection != none {
		if self.rotationDirection == right {
			self.rotationAngle = math.Mod(self.rotationAngle+angleSpeed*delta, math.Pi)
		} else {
			self.rotationAngle = math.Mod(self.rotationAngle-angleSpeed*delta, math.Pi)
		}
	}
	self.speedRatio = lerp(self.speedRatio, self.targetSpeedRatio, delta)
	self.Location.X += (self.SelectedShip.Speed * self.speedRatio * delta) * math.Cos(self.rotationAngle)
	self.Location.Y += (self.SelectedShip.Speed * self.speedRatio * delta) * math.Sin(self.rotationAngle)
	self.Unlock()
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}
