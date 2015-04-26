package commonStructs

type Ship struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Weight      int     `json:"weight"`
	CannonCount byte    `json:"cannonCount"`
	Speed       float64 `json:"speed"`
	HP          int     `json:"baseHP"`
	CurrentHP   int     `json:"currentHP"`
}
