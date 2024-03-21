package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	username string
	conn     *websocket.Conn
}

func NewClient(username string, conn *websocket.Conn) *Client {
	return &Client{username, conn}
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient("", conn)

	log.Println("New client has joined: ", client)
}
