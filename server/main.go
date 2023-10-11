package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func main() {
	http.HandleFunc("/chat", handleConnections)
	go handleMessages()
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println(len(clients))
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg string 
		_, p, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			return
		}
		msg = string(p[:])
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
