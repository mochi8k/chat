package main

import (
	"github.com/gorilla/websocket"
)

// one person
type client struct {
	socket *websocket.Conn

	// send message channel.
	send chan []byte

	// chat room
	room *room
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name.c.userData["name"].(string)
			c.room.forward <-msgLLL
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
