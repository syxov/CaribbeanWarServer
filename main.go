package main

import (
	"CaribbeanWarServer/api"
	"CaribbeanWarServer/world"
	"net/http"
	"os"
	"runtime"
)

var (
	worldStr world.WorldStruct
	dbConn   api.DbConnection
)

func init() {
	dbConn.Open()
	worldStr = world.WorldStruct{DbConn: &dbConn}
	runtime.GOMAXPROCS(runtime.NumCPU())
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "80")
	}
}

func main() {
	defer dbConn.Close()
	http.HandleFunc("/ws", api.Handler(&worldStr, dbConn))
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
