package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	gameManager *GameManager
	idx         int
}

func NewServer() *Server {
	return &Server{
		gameManager: NewGameManager(),
		idx:         0,
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	id := s.idx
	s.idx++

	s.gameManager.AddPlayer(id, conn)

	s.messageLoop(id, conn)
}

func (s *Server) messageLoop(id int, conn *websocket.Conn) {
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", data)

		s.gameManager.HandleMessage(id, messageType, data)

		// if err := conn.WriteMessage(messageType, data); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
}
