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
	angleSpeed = 10
	velocity   = float64(1)
)

type User struct {
	ID                uint
	Email             string
	Nick              string
	Cash              uint
	Conn              *websocket.Conn
	Location          *Point
	Ships             []Ship
	SelectedShip      *Ship
	NearestUsers      []*User
	InWorld           bool
	targetSpeedRatio  float64
	speedRatio        float64
	rotationAngle     float64
	rotationDirection byte
	sync.RWMutex
}

func (self *User) Bounds() *rtree.Rect {
	bound, _ := rtree.NewRect(rtree.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	return bound
}

func (self *User) SetMove(moveType string) {
	self.Lock()
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
	self.Unlock()
}

func (self *User) UpdatePosition(realDelta float64) {
	self.Lock()
	if self.rotationDirection != none {
		if self.rotationDirection == right {
			self.rotationAngle = math.Mod(self.rotationAngle+angleSpeed*realDelta, math.Pi)
		} else {
			self.rotationAngle = math.Mod(self.rotationAngle-angleSpeed*realDelta, math.Pi)
		}
	}
	self.Location.X += (self.SelectedShip.Speed * self.speedRatio * realDelta) * math.Cos(self.rotationAngle)
	self.Location.Y += (self.SelectedShip.Speed * self.speedRatio * realDelta) * math.Sin(self.rotationAngle)
	self.speedRatio = lerp(self.speedRatio, self.targetSpeedRatio, 0.01*realDelta)
	self.Unlock()
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}
