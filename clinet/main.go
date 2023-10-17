package main

import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
	"time"
)

func main() {
	serverAddr := "ws://localhost:8080/chat"

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Connection closed:", err)
				return
			}
			fmt.Println("Received:", string(msg))
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		var msg string
		fmt.Scan(&msg)
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
