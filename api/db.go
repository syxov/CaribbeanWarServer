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
		shipId   uint
		cash     uint
		nick     string
		rotation float64
		location []uint8

		name        string
		weight      uint16
		cannonCount byte
		speed       float64
		hp          uint16
	)
	err := self.db.QueryRow(`
		SELECT id, cash, nick, location, rotation FROM users 
		WHERE email=$1 AND password=$2
	`, email, password).Scan(&id, &cash, &nick, &location, &rotation)
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
		rows.Scan(&shipId, &name, &weight, &cannonCount, &speed, &hp)
		ships = append(ships, structs.Ship{
			ID:          shipId,
			Name:        name,
			Weight:      weight,
			CannonCount: cannonCount,
			Speed:       speed,
			HP:          hp,
		})
	}
	return &structs.User{
		ID:            id,
		Email:         email,
		Cash:          cash,
		Nick:          nick,
		Location:      &structs.Point{float64(location[0]), float64(location[1])},
		Ships:         ships,
		RotationAngle: rotation,
	}, nil
}

func (self *DbConnection) SaveUserLocation(user *structs.User) error {
	_, err := self.db.Query(`
		UPDATE users
		SET location=$1, rotation=$2
		WHERE id=$3
	`, user.Location, user.RotationAngle, user.ID)
	return err
}
