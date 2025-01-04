package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type GameManager struct {
	players map[int]*Player
}

type Player struct {
	conn *websocket.Conn
	mu   *sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		players: make(map[int]*Player),
	}
}

func (gm *GameManager) AddPlayer(id int, conn *websocket.Conn) {
	gm.Broadcast(websocket.TextMessage, []byte(fmt.Sprintf("New player: %d", id)))
	gm.players[id] = &Player{conn: conn, mu: &sync.Mutex{}}
	gm.SendMsg(id, websocket.TextMessage, []byte(fmt.Sprintf("Init %d %d", id, len(gm.players))))
}

func (gm *GameManager) RemovePlayer(id int) {
	gm.Broadcast(websocket.TextMessage, []byte(fmt.Sprintf("Player left: %d", id)))
	delete(gm.players, id)
}

func (gm *GameManager) Broadcast(messageType int, data []byte) {
	for _, player := range gm.players {
		player.mu.Lock()
		if err := player.conn.WriteMessage(messageType, data); err != nil {
			return
		}
		player.mu.Unlock()
	}
}

func (gm *GameManager) HandleMessage(id int, messageType int, data []byte) {
	gm.Broadcast(messageType, data)
}

func (gm *GameManager) Close() {
	for _, player := range gm.players {
		player.conn.Close()
	}
}

func (gm *GameManager) SendMsg(id int, messageType int, data []byte) {
	player := gm.players[id]
	player.mu.Lock()
	defer player.mu.Unlock()

	if err := player.conn.WriteMessage(messageType, data); err != nil {
		return
	}
}
