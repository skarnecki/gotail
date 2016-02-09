package pump

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

//WebHandler handles http/ws
type WebHandler struct {
	Filechannel chan string
	Buffer      []string
	BufferSize  int
}

func (wh *WebHandler) Websocket(ws *websocket.Conn) {
	for _, line := range wh.Buffer {
		ws.Write([]byte(line))
	}
	for {
		select {
		case msg := <-wh.Filechannel:
			fmt.Printf("new data")
			ws.Write([]byte(msg))
			wh.Buffer = append(wh.Buffer, msg)
			wh.Buffer = wh.Buffer[len(wh.Buffer)-wh.BufferSize:]
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
