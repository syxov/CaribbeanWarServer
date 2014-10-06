package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

const dbUrl = "postgres://wxkjimhylvjmxu:SYNQgXeUGnWQQbge7LwyPka3SB@ec2-54-225-255-208.compute-1.amazonaws.com:5432/delbcum37jd7n5"

var DbConn = new(dbConnection)

type dbConnection struct {
	db *sql.DB
}

func (self *dbConnection) Close() {
	self.db.Close()
}

func (self *dbConnection) Open() {
	db, err := sql.Open("postgres", dbUrl)
	self.db = db
	if err != nil {
		log.Fatal(err)
	}
}

func (self *dbConnection) CheckUserExist(email, password string) bool {
	if rows, err := self.db.Query("SELECT id FROM users WHERE email='?' AND password='?'", email, password); err == nil {
		return rows.Next()
	} else {
		log.Fatal(err)
	}
	return false
}

type userInfo struct {
	ID    uint
	Email string
	Cash uint
}

func (self *dbConnection) GetUserInfo(email, password string) *userInfo {
	var (
		id uint
		cash uint
	)
	err := self.db.QueryRow("SELECT id, cash FROM users WHERE email='" + email + "'AND password='" + password + "'").Scan(&id, &cash)
	if err != nil {
		log.Print(err)
		return nil
	}
	result := userInfo{id, email, cash}
	return &result
}
