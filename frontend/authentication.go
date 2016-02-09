package frontend

import (
	"crypto/sha1"
	"encoding/base64"
)

//Authentication static content handler
type Authentication struct {
	User     string
	Password string
}

//Secret password validation
func (sh *Authentication) Secret(incomingUser, realm string) string {
	if incomingUser == sh.User {
		d := sha1.New()
		d.Write([]byte(sh.Password))
		return "{SHA}" + base64.StdEncoding.EncodeToString(d.Sum(nil))
	}
	return ""
}
