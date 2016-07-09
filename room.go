package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte

	// for joiner
	join chan *client

	// for leavener
	leave chan *client

	clients map[*client]bool
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHttp(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHttp:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client

	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true

		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)

		case msg := <-r.forward:
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// 送信
				default:
					// 失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
