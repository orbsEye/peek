package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"peek/api/handlers/coincap"
	"peek/api/handlers/hackernews"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	assets := []string{"bitcoin", "ethereum", "litecoin", "monero"}
	coincap.StartWebSocketServer(assets, clientConn)
}

func main() {
	http.HandleFunc("/hackernews", hackernews.HackerNewsHandler)

	http.HandleFunc("/coincap", websocketHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	port := ":8080"
	fmt.Printf("Starting server on %s\n", port)

	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	fmt.Println("Server shutting down...")
}
