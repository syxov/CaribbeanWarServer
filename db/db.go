package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Open() {
	db, _ = sql.Open("postgres", "postgres://czegnjaplurnsd:OAuYt0yTCAdyX5BpcAcMRYjg7U@ec2-54-225-101-202.compute-1.amazonaws.com:5432/dbv0q37nim28dm")
}

func CheckUserExist(email, password string) bool {
	rows, _ := db.Query("SELECT id FROM users WHERE email='" + email + "' AND password='" + password + "'")
	return rows.Next()
}
