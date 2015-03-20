package messagesStructs

import (
	"CaribbeanWarServer/point"
)

type MoveIncome struct {
	Details struct {
		Type string `json:"type"`
	} `json:"details"`
	Message
}

type MoveOutcome struct {
	Details MoveOutcomeDetails `json:"details"`
	Action  string             `json:"action"`
}

type MoveOutcomeDetails struct {
	ID       uint          `json:"id"`
	Type     string        `json:"type"`
	Alpha    float64       `json:"alpha"`
	Location point.Point2D `json:"location"`
}
