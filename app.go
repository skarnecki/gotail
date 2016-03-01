package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"github.com/skarnecki/gotail/frontend"
	"github.com/skarnecki/gotail/pump"
	"golang.org/x/net/websocket"
)

var (
	filename = kingpin.Arg("filename", "Path to tailed file.").Required().ExistingFile()
	number   = kingpin.Flag("number", "Starting lines number.").Default("10").Int()
	host     = kingpin.Flag("host", "Listening host, default 0.0.0.0").Default("0.0.0.0").IP()
	port     = kingpin.Flag("port", "listening port, default 9001").Default("9001").Int()
	cert     = kingpin.Flag("cert", "path to cert file (HTTPS)").ExistingFile()
	key      = kingpin.Flag("key", "path to key file (HTTPS)").ExistingFile()
	user     = kingpin.Flag("user", "Basic auth user").String()
	password = kingpin.Flag("password", "Basic auth password").String()
)

func main() {
	kingpin.Parse()
	filechannel := make(chan string, 100)
	go pump.TailFile(filechannel, *filename)

	//todo forbid basicauth w/o user & pass
	mainpage := &frontend.MainPage{HTTPSMode: false, BasicAuth: false, UserName: *user, UserPassword: *password}

	handler := pump.WebHandler{Filechannel: filechannel, Buffer: make([]string, *number), BufferSize: *number}
	address := fmt.Sprintf("%s:%d", *host, *port)

	r := mux.NewRouter()
	r.Handle("/", mainpage)
	r.Handle("/socket", websocket.Handler(handler.Websocket))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", r)
	fmt.Printf("Listening on %s\n", address)
	if *cert != "" && *key != "" {
		http.ListenAndServeTLS(address, *cert, *key, nil)
	}
	http.ListenAndServe(address, nil)
}
