package Ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{} // Mutex für Thread-Safety

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Client hinzufügen
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Client entfernen beim Beenden
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received from client: %s", msg)

		// Echo direkt an alle
		BroadcastMessage(msg)
	}
}

func BroadcastMessage(msg []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			err2 := client.Close()
			if err2 != nil {
				return
			}
			delete(clients, client)
		}
	}
}
