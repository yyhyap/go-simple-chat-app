package main

import (
	"flag"
	"fmt"
	"go-simple-chat-app/server/chat"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Hello")
	chat.Start(*port)
}
