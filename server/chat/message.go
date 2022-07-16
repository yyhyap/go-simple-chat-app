package chat

import "go-simple-chat-app/server/utils"

type Message struct {
	ID     int64  `json:"id"`
	Body   string `json:"body"`
	Sender string `json:"sender"`
}

func NewMessage(body, sender string) *Message {
	return &Message{
		ID:     utils.GetRandomInt64(),
		Body:   body,
		Sender: sender,
	}
}
