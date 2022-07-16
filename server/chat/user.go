package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type User struct {
	Username string
	Conn     *websocket.Conn
	Global   *Chat
}

func (u *User) Read() {
	// listening for message
	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			log.Println("Error occured when reading message: ", err.Error())
			break
		} else {
			u.Global.Messages <- NewMessage(string(message), u.Username)
		}
	}

	u.Global.Leave <- u
}

func (u *User) Write(message *Message) {
	bytes, err := json.Marshal(message)

	if err != nil {
		log.Println("Error occured when marshalling message: ", err.Error())
	}

	if err := u.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
		log.Println("Error occured when writing message: ", err.Error())
	}
}
