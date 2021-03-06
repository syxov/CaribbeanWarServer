package messagesStructs

import (
	"CaribbeanWarServer/point"
)

type ShootIncome struct {
	Details struct {
		Location  point.Point3D `json:"location"`
		Angle     float64       `json:"angle"`
		Direction float64       `json:"direction"`
	} `json:"details"`
	Message
}

type ShootOutcome struct {
	Details ShootOutcomeDetails `json:"details"`
	Action  string              `json:"action"`
}

type ShootOutcomeDetails struct {
	ID        uint          `json:"id"`
	Alpha     float64       `json:"alpha"`
	Angle     float64       `json:"angle"`
	Location  point.Point3D `json:"location"`
	Direction float64       `json:"direction"`
}

type Hit struct {
	Action  string     `json:"action"`
	Details HitDetails `json:"details"`
}

type HitDetails struct {
	ID       uint
	Location point.Point3D
	Damage   uint
}
