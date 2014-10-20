package structs

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID           uint
	Email        string
	Nick         string
	Cash         uint
	Conn         *websocket.Conn
	Location     Point
	Ships        []Ship
	SelectedShip *Ship
}
