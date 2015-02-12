package structs

type Ship struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Weight      uint16  `json:"weight"`
	CannonCount byte    `json:"cannonCount"`
	Speed       float64 `json:"speed"`
	HP          uint16  `json:"hp"`
}
