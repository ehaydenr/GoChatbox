package main

import (
	"github.com/ehaydenr/chatbox/chat"
	"html/template"
	"net/http"
	"os"
)

var client_id string = func() string {
	if s := os.Getenv("client_id"); s != "" {
		return s
	}
	return "1017771976315-g6bc9dc2a6ud3v4ngare0fslpgf4lmqb.apps.googleusercontent.com"
}()

var chat_template *template.Template = template.Must(template.ParseFiles("templates/chat.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	chat_template.Execute(w, client_id)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	server := chat.NewServer("/ws")
	go server.Listen()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
