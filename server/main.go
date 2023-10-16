package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
type function func([]string, gorm.DB)(bool)

var routes = map[string]function{
	"create":createBulshit,
}


var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []string)






/// main function and connections
func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Bulshit{})
	http.HandleFunc("/chat", handleConnections)
	// go handleMessages()
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	clients[ws] = true
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	for {
		var msg []string
		_, p, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			return
		}
		msg = strings.Split(string(p[:]), ":::")
		ans := routes[msg[0]](msg, *db)
		err = ws.WriteMessage(websocket.TextMessage, []byte(strconv.FormatBool(ans)))
		if err!=nil{
			ws.Close()
			delete(clients, ws)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg[0]))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}


/// models 

type Bulshit struct {
	gorm.Model
	Name string
	Code   string
	CreatedAt time.Time  
}





/// apis 
func createBulshit(request []string, db gorm.DB)bool{
	if len(request) < 3 {
		return false
	}
	name := request[1]
	code := request[2]
	db.Create(&Bulshit{Name: name, Code: code, CreatedAt: time.Now()})
	return true
}

func ss(){
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Bulshit{})

	db.Create(&Bulshit{Name: "ss", Code: "100"})

	var user Bulshit
	// db.First(&user, 1)                 // find product with integer primary key
	db.First(&user, "Name = ?", "ss") // find product with code D42

	db.Model(&user).Update("code", 200)
	db.Model(&user).Updates(Bulshit{Name: "hello"})

	fmt.Println(user.Code)

}
