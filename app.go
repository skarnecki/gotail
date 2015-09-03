package main

import (
	"net/http"
	"flag"
	"os"
	"log"
	"os/exec"
	"golang.org/x/net/websocket"
	"time"
	"fmt"
)

func tailFile(filename string, cs chan string) {
	cmd := exec.Command("tail", "-f", filename)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()

	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	for {
		readTo := make([]byte, 1024)
		stdout.Read(readTo)
		log.Print("fileread")
		cs <- string(readTo)
	}
}

func main() {
	filename := flag.String("filename", "", "path to tailed file")
	flag.Parse()

	if _, err := os.Stat(*filename);

	os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", filename)
		return
	}

	filechannel := make(chan string, 100)
	go tailFile(*filename, filechannel)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/socket", websocket.Handler(func (ws *websocket.Conn) {
		for {
			select {
			case msg := <-filechannel:
				log.Print("sending")
				ws.Write([]byte(msg))
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}))
	http.ListenAndServe(":8080", nil)
}