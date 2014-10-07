package main

import (
	"CaribbeanWarServer/api"
	"net/http"
	"os"
	"runtime"
)

func init() {
	api.DbConn.Open()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set by heroku
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "80")
	}
}

func main() {
	defer func() {
		api.DbConn.Close()
	}()
	http.HandleFunc("/ws", api.Handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
