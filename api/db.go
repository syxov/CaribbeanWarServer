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
		id                                 uint
		shipId                             uint
		cash                               uint
		nick                               string
		rotation, coordinateX, coordinateY float64

		name        string
		weight      uint16
		cannonCount byte
		speed       float64
		hp          uint16
	)
	err := self.db.QueryRow(`
		SELECT id, cash, nick, coordinate_x, coordinate_y, rotation FROM users 
		WHERE email=$1 AND password=$2
	`, email, password).Scan(&id, &cash, &nick, &coordinateX, &coordinateY, &rotation)
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
		Location:      &structs.Point{coordinateX, coordinateY},
		Ships:         ships,
		RotationAngle: rotation,
	}, nil
}

func (self *DbConnection) SaveUserLocation(user *structs.User) error {
	_, err := self.db.Query(`
		UPDATE users
		SET coordinate_x=$1, coordinate_y=$2, rotation=$3
		WHERE id=$4
	`, user.Location.X, user.Location.Y, user.RotationAngle, user.ID)
	return err
}
