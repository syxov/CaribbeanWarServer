package main

import (
	"CaribbeanWarServer/api"
	"CaribbeanWarServer/harbor"
	"net/http"
	"os"
	"runtime"
)

var (
	dbConn    api.DbConnection
	harborStr harbor.HarborStruct
)

func init() {
	dbConn.Open()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "80")
	}
}

func main() {
	defer dbConn.Close()
	http.HandleFunc("/ws", api.Handler(&harborStr, dbConn))
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
