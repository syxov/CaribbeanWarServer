package main

import (
	"CaribbeanWarServer/db"
	"CaribbeanWarServer/socket"
	"net/http"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	db.Open()
}

func main() {
	http.HandleFunc("/ws", socket.Handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
