package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"

	auth "github.com/abbot/go-http-auth"
	"github.com/alecthomas/kingpin"
	"github.com/skarnecki/gotail/pump"
	"golang.org/x/net/websocket"
)

//StaticHandler static content handler
type StaticHandler struct {
	User     string
	Password string
}

//Secret password validation
func (sh *StaticHandler) Secret(incomingUser, realm string) string {
	if incomingUser == sh.User {
		d := sha1.New()
		d.Write([]byte(sh.Password))
		return "{SHA}" + base64.StdEncoding.EncodeToString(d.Sum(nil))
	}
	return ""
}

func (sh *StaticHandler) handle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	http.FileServer(http.Dir("./static")).ServeHTTP(w, &r.Request)
}

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

	http.Handle("/", initContentHandler(*user, *password))
	handler := pump.WebHandler{Filechannel: filechannel, Buffer: make([]string, *number), BufferSize: *number}
	http.Handle("/socket", websocket.Handler(handler.Websocket))

	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Listening on %s\n", address)
	if *cert != "" && *key != "" {
		http.ListenAndServeTLS(address, *cert, *key, nil)
	}
	http.ListenAndServe(address, nil)
}

func initContentHandler(user, password string) http.Handler {
	if user != "" && password != "" {
		staticContentHandler := &StaticHandler{User: user, Password: password}
		authWrapper := auth.NewBasicAuthenticator("Gotail", staticContentHandler.Secret)
		return authWrapper.Wrap(staticContentHandler.handle)
	}
	return http.FileServer(http.Dir("./static"))
}
