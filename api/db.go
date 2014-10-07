package api

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const dbUrl = "postgres://wxkjimhylvjmxu:SYNQgXeUGnWQQbge7LwyPka3SB@ec2-54-225-255-208.compute-1.amazonaws.com:5432/delbcum37jd7n5"

var DbConn = new(dbConnection)

type dbConnection struct {
	db *sql.DB
}

func (self *dbConnection) Open() {
	self.db, _ = sql.Open("postgres", dbUrl)
}

func (self *dbConnection) Close() {
	self.db.Close()
}

func (self *dbConnection) CheckUserExist(email, password string) bool {
	if rows, err := self.db.Query("SELECT id FROM users WHERE email=$1 AND password=$2", email, password); err == nil {
		return rows.Next()
	} else {
		return false
	}
}

type UserInfo struct {
	ID    uint
	Email string
	Cash  uint
}

func (self *dbConnection) GetUserInfo(email, password string) *UserInfo {
	var (
		id   uint
		cash uint
	)
	err := self.db.QueryRow("SELECT id, cash FROM users WHERE email=$1 AND password=$2", email, password).Scan(&id, &cash)
	if err != nil {
		return nil
	}
	result := UserInfo{id, email, cash}
	return &result
}
