package api

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DbConn = new(dbConnection)

type dbConnection struct {
	db *sql.DB
}

func (self *dbConnection) Open() {
	self.db, _ = sql.Open("postgres", "postgres://czegnjaplurnsd:OAuYt0yTCAdyX5BpcAcMRYjg7U@ec2-54-225-101-202.compute-1.amazonaws.com:5432/dbv0q37nim28dm")
}

func (self *dbConnection) CheckUserExist(email, password string) bool {
	stm, _ := self.db.Prepare("SELECT id FROM users WHERE email='" + email + "' AND password='" + password + "'")
	rows, _ := stm.Query()
	return rows.Next()
}
