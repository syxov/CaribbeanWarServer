package structs

import (
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/rtree"
	"math"
)

type Core struct {
	StartPosition, CurrentPosition, previousPosition point.Point3D
	ShipRotationAngle, ShootAngle                    float64
	ID                                               uint
	time                                             float64
	coefX, coefY, coefZ                              float64
}

func NewCore(position point.Point3D, shipRotationAngle, shootAngle, direction float64, ID uint) *Core {
	core := &Core{
		StartPosition:     position,
		CurrentPosition:   position,
		previousPosition:  position,
		ShipRotationAngle: shipRotationAngle - direction*math.Pi/2.0,
		ShootAngle:        shootAngle,
		ID:                ID,
		time:              0,
	}
	core.coefX = math.Cos(core.ShootAngle) * math.Cos(math.Pi-core.ShipRotationAngle) * speed
	core.coefY = math.Cos(core.ShootAngle) * math.Sin(math.Pi-core.ShipRotationAngle) * speed
	core.coefZ = math.Sin(core.ShootAngle) * speed
	return core
}

const speed float64 = 100

func (self *Core) UpdatePosition(delta float64) {
	self.time += delta
	self.previousPosition = self.CurrentPosition
	self.CurrentPosition = point.Point3D{
		X: self.StartPosition.X + self.time*self.coefX,
		Y: self.StartPosition.Y + self.time*self.coefY,
		Z: self.StartPosition.Z + self.time*self.coefZ - 9.8*math.Pow(self.time, 2.0)/2.0,
	}
}

func (self *Core) OverWater() bool {
	return self.CurrentPosition.Z > 0
}

const radius float64 = 1

func (self *Core) GetBounds(r ...int) *rtree.Rect {
	lenBetweenPoint := math.Sqrt(math.Pow(self.CurrentPosition.X-self.previousPosition.X, 2) + math.Pow(self.CurrentPosition.Y-self.previousPosition.Y, 2))
	return rtree.NewRect(rtree.Point{self.previousPosition.X - radius/2.0, self.previousPosition.Y - radius/2.0}, []float64{radius + lenBetweenPoint, radius}, self.ShipRotationAngle)
}
