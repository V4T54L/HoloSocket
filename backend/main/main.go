package main

import (
	"fmt"
	"net/http"
)

func main() {
	s := NewServer()

	http.HandleFunc("/ws", s.handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	fmt.Println("WebSocket server is running on :8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
