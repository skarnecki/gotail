package frontend

import (
	"net/http"
	"text/template"

	httpauth "github.com/abbot/go-http-auth"
)

const (
	WebsocketProtocol       = "ws"
	WebsocketSecureProtocol = "wss"
)

type MainPage struct {
	HTTPSMode    bool
	BasicAuth    bool
	UserName     string
	UserPassword string
}

type MainPageDetails struct {
	Title      string
	WSProtocol string
}

func (mp *MainPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mp.BasicAuth {
		staticContentHandler := &Authentication{User: mp.UserName, Password: mp.UserPassword}
		authWrapper := httpauth.NewBasicAuthenticator("Gotail", staticContentHandler.Secret)
		authWrapper.Wrap(mp.AuthTail).ServeHTTP(w, r)
	} else {
		mp.Tail(w, r)
	}
}

func (mp *MainPage) AuthTail(w http.ResponseWriter, r *httpauth.AuthenticatedRequest) {
	mp.Tail(w, &r.Request)
}

func (mp *MainPage) Tail(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	data := &MainPageDetails{Title: "Gotail", WSProtocol: WebsocketProtocol}
	if mp.HTTPSMode {
		data.WSProtocol = WebsocketSecureProtocol
	}
	t.Execute(w, data)
}
