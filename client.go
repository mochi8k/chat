package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// one person
type client struct {
	socket *websocket.Conn

	// send message channel.
	send chan *message

	// chat room
	room *room

	// ユーザー情報
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			fmt.Println(msg.AvatarURL)
			c.room.forward <- msg
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
