package chat

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go-simple-chat-app/server/utils"

	"github.com/gorilla/websocket"
)

type Chat struct {
	Users    map[string]*User
	Messages chan *Message
	Join     chan *User
	Leave    chan *User
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("%s %s%s %v\n", r.Method, r.Host, r.URL, r.Proto)
		return r.Method == http.MethodGet
	},
}

func (c *Chat) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Error occured on websocket connection: ", err.Error())
		return
	}

	keys := r.URL.Query()
	username := keys.Get("username")

	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("AnonymousUser-%d", utils.GetRandomInt64())
	}

	user := &User{
		Username: username,
		Conn:     conn,
		Global:   c,
	}

	c.Join <- user

	user.Read()
}

func (c *Chat) Run() {
	for {
		select {
		case user := <-c.Join:
			c.add(user)
		case message := <-c.Messages:
			c.broadcast(message)
		case user := <-c.Leave:
			c.disconnect(user)
		}
	}
}

func (c *Chat) add(user *User) {
	if _, exist := c.Users[user.Username]; !exist {
		c.Users[user.Username] = user

		body := fmt.Sprintf("%s has joined the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))

		log.Printf("Added user: %s, Total users: %d\n", user.Username, len(c.Users))
	}
}

func (c *Chat) broadcast(message *Message) {
	log.Printf("Broadcasting message : %v\n", message)
	for _, user := range c.Users {
		user.Write(message)
	}
}

func (c *Chat) disconnect(user *User) {
	if _, exist := c.Users[user.Username]; exist {
		defer user.Conn.Close()
		delete(c.Users, user.Username)

		body := fmt.Sprintf("%s has left the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))

		log.Printf("User %s has left the chat, Total users: %d\n", user.Username, len(c.Users))
	}
}

func Start(port string) {
	log.Printf("Chat listening on http://localhost:%s\n", port)

	c := &Chat{
		Users:    make(map[string]*User),
		Messages: make(chan *Message),
		Join:     make(chan *User),
		Leave:    make(chan *User),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go webchat"))
	})

	http.HandleFunc("/chat", c.Handler)

	go c.Run()

	log.Fatal(http.ListenAndServe(port, nil))
}
