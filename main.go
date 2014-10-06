package main

import (
	"CaribbeanWarServer/api"
	"net/http"
	"net/smtp"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	api.DbConn.Open()
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "80")
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			auth := smtp.PlainAuth("", "al.syxov@gmail.com", "505604qw", "smtp.gmail.com")
			smtp.SendMail("smtp.gmail.com:587", auth, "al.syxov@gmail.com", []string{"al.syxov@gmail.com"}, []byte("try to send"))
		}
	}()
	defer func() {
		api.DbConn.Close()
	}()
	http.HandleFunc("/ws", api.Handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
