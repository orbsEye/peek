package main

import (
	"net/http"
	"os"
	"os/signal"
	"peek/api/handlers/coincap"
	"peek/api/handlers/hackernews"
)

func main() {

	http.HandleFunc("/hackernews", hackernews.HackerNewsHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := ":8080"
	println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}

	coincap.CoincapInfo()
	assets := []string{"bitcoin", "ethereum", "litecoin", "monero"}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go coincap.StartWebSocket(assets)
	<-interrupt
}
