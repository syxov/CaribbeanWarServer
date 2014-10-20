package api

import (
	"CaribbeanWarServer/structs"
	"database/sql"
	_ "github.com/lib/pq"
)

const dbUrl = "postgres://wxkjimhylvjmxu:SYNQgXeUGnWQQbge7LwyPka3SB@ec2-54-225-255-208.compute-1.amazonaws.com:5432/delbcum37jd7n5"

type DbConnection struct {
	db *sql.DB
}

func (self *DbConnection) Open() {
	self.db, _ = sql.Open("postgres", dbUrl)
}

func (self *DbConnection) Close() {
	self.db.Close()
}

func (self *DbConnection) GetUserInfo(email, password string) (*structs.User, error) {
	var (
		id       uint
		cash     uint
		nick     string
		location []uint8

		name        string
		weight      uint16
		cannonCount byte
		speed       uint16
		hp          uint16
	)
	err := self.db.QueryRow(`
		SELECT id, cash, nick, location FROM users 
		WHERE email=$1 AND password=$2
	`, email, password).Scan(&id, &cash, &nick, &location)
	if err != nil {
		return nil, err
	}
	rows, err := self.db.Query(`
		SELECT * FROM ships
		WHERE id IN (
			SELECT ship_id FROM user_ships
			WHERE user_id=$1
		)
	`, id)
	if err != nil {
		return nil, err
	}
	ships := []structs.Ship{}
	for rows.Next() {
		rows.Scan(&id, &name, &weight, &cannonCount, &speed, &hp)
		ships = append(ships, structs.Ship{
			ID:          id,
			Name:        name,
			Weight:      weight,
			CannonCount: cannonCount,
			Speed:       speed,
			HP:          hp,
		})
	}
	return &structs.User{
		ID:           id,
		Email:        email,
		Cash:         cash,
		Nick:         nick,
		Location:     structs.Point{float64(location[0]), float64(location[1])},
		Ships:        ships,
		SelectedShip: nil,
	}, nil
}
