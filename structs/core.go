package structs

import (
	"CaribbeanWarServer/rtree"
	"math"
)

type Core struct {
	StartPosition   Point3D
	Angle, Alpha    float64
	ID              uint
	CurrentPosition Point3D
	time            float64
}

func NewCore(position *Point3D, angle, direction float64, ID uint) *Core {
	return &Core{
		StartPosition:   *position,
		CurrentPosition: *position,
		Angle:           angle,
		Alpha:           -angle - direction*math.Pi/2,
		ID:              ID,
		time:            0,
	}
}

const speed = 100

func (self *Core) UpdatePosition(delta float64) {
	self.time += delta
	self.CurrentPosition = Point3D{
		X: self.StartPosition.X + speed*self.time*math.Cos(self.Angle)*math.Cos(self.Alpha),
		Y: self.StartPosition.Y + speed*self.time*math.Sin(self.Angle) - 9.8*math.Pow(self.time, 2)/2,
		Z: self.StartPosition.Z + speed*self.time*math.Cos(self.Angle)*math.Sin(self.Alpha),
	}
}

func (self *Core) UnderWater() bool {
	return self.CurrentPosition.Z <= 0
}

const radius = 1

func (self *Core) GetBounds(r ...int) *rtree.Rect {
	rect, _ := rtree.NewRect(rtree.Point{self.CurrentPosition.X - radius/2, self.CurrentPosition.Y - radius/2}, []float64{radius, radius})
	return rect
}
