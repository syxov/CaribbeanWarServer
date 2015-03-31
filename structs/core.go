package structs

import (
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/rtree"
	"math"
)

type Core struct {
	StartPosition, CurrentPosition, previousPosition point.Point3D
	Angle, Alpha                                     float64
	ID                                               uint
	time                                             float64
}

func NewCore(position point.Point3D, alpha, angle, direction float64, ID uint) *Core {
	return &Core{
		StartPosition:    position,
		CurrentPosition:  position,
		previousPosition: position,
		Angle:            angle - direction*math.Pi/2.0,
		Alpha:            alpha,
		ID:               ID,
		time:             0,
	}
}

const speed float64 = 100

func (self *Core) UpdatePosition(delta float64) {
	self.time += delta
	self.previousPosition = self.CurrentPosition
	self.CurrentPosition = point.Point3D{
		X: self.StartPosition.X + speed*self.time*math.Cos(self.Alpha)*math.Cos(self.Angle),
		Y: self.StartPosition.Y + speed*self.time*math.Sin(self.Alpha) - 9.8*math.Pow(self.time, 2)/2.0,
		Z: self.StartPosition.Z + speed*self.time*math.Cos(self.Alpha)*math.Sin(self.Angle),
	}
}

func (self *Core) UnderWater() bool {
	return self.CurrentPosition.Y <= 0
}

const radius float64 = 1

func (self *Core) GetBounds(r ...int) *rtree.Rect {
	lenBetweenPoint := math.Sqrt(math.Pow(self.CurrentPosition.X-self.previousPosition.X, 2) + math.Pow(self.CurrentPosition.Z-self.previousPosition.Z, 2))
	return rtree.NewRect(rtree.Point{self.previousPosition.X - radius/2.0, self.previousPosition.Z - radius/2.0}, []float64{radius + lenBetweenPoint, radius}, self.Angle)
}
