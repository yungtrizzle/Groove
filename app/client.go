package app

import (
	"bytes"
	"github.com/gorilla/websocket"
)

type Client struct {
	User   string
	Chatid int
	Room   string //curent room
	Roomid int

	send chan []byte

	conn *websocket.Conn
}
