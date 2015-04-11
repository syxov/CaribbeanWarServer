package messagesStructs

import (
	"CaribbeanWarServer/point"
)

type RespawnIncome Message

type RespawnOutcome struct {
	Message
	Details RespawnOutcomeDetails `json:"details"`
}

type RespawnOutcomeDetails struct {
	Location point.Point2D `json:"position"`
	Rotation float64       `json:"rotation"`
}