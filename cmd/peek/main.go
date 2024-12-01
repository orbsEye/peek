package main

import (
	"os"
	"os/signal"
	"peek/api/handlers/coincap"
	"peek/api/handlers/hackernews"
)

func main() {

	hackernews.HackerNewsInfo()

	coincap.CoincapInfo()
	assets := []string{"bitcoin", "ethereum", "litecoin", "monero"}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go coincap.StartWebSocket(assets)
	<-interrupt
}
