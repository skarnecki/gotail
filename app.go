package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
)

//WebHandler handles http/ws
type WebHandler struct {
	Filechannel chan string
	Buffer      []string
	BufferSize  int
}

func (wh *WebHandler) websocketPump(ws *websocket.Conn) {
	for _, line := range wh.Buffer {
		ws.Write([]byte(line))
	}
	for {
		select {
		case msg := <-wh.Filechannel:
			ws.Write([]byte(msg))
			wh.Buffer = append(wh.Buffer, msg)
			wh.Buffer = wh.Buffer[len(wh.Buffer)-wh.BufferSize:]
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

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

func main() {
	filename := flag.String("filename", "", "path to tailed file")
	number := flag.Int("number", 10, "starting lines number, default 10")
	host := flag.String("host", "0.0.0.0", "listening host, default 0.0.0.0")
	port := flag.Int("port", 9001, "listening port, default 9001")
	cert := flag.String("cert", "", "path to cert file (HTTPS)")
	key := flag.String("key", "", "path to key file (HTTPS)")

	flag.Parse()
	if *filename == "" {
		fmt.Println("Specify tailed file with --filename")
		return
	}

	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		fmt.Printf("Can't read file: %s\n", *filename)
		return
	}

	filechannel := make(chan string, 100)
	go func(cs chan string) {
		// file, _ := tail.TailFile(*filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: -1000, Whence: os.SEEK_END}})
		file, _ := tail.TailFile(*filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_CUR}})
		for line := range file.Lines {
			cs <- string(line.Text)
		}
	}(filechannel)

	staticContentHandler := &StaticHandler{User: "john", Password: "cena"}
	authWrapper := auth.NewBasicAuthenticator("example.com", staticContentHandler.Secret)
	http.Handle("/", authWrapper.Wrap(staticContentHandler.handle))
	handler := WebHandler{Filechannel: filechannel, Buffer: make([]string, *number), BufferSize: *number}
	http.Handle("/socket", websocket.Handler(handler.websocketPump))

	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Listening on %s", address)
	if *cert != "" && *key != "" {
		http.ListenAndServeTLS(address, *cert, *key, nil)
	}
	http.ListenAndServe(address, nil)
}
