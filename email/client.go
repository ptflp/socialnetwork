package email

import (
	"errors"
	"log"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func Send() {
	// Choose auth method and set it up
	auth := LoginAuth("user", "pass")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"to@example.com"}
	msg := []byte("To: to@example.com\r\n" +
		"Subject: New Hack\r\n" +
		"\r\n" +
		"Wonderful solution\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "from@example.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
