package main

import (
	"CaribbeanWarServer/api"
	"CaribbeanWarServer/auth"
	"CaribbeanWarServer/harbor"
	//	"CaribbeanWarServer/world"
	"net/http"
	"os"
	"runtime"
)

var dbConn api.DbConnection

func init() {
	dbConn.Open()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set by heroku
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "80")
	}
}

func main() {
	defer dbConn.Close()
	http.HandleFunc("/ws", auth.Handler(&dbConn, harbor.GetHarbor()))
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
