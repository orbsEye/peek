package coincap

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Coincap struct {
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price_usd"`
}

func StartWebSocket(assets []string) {
	url := fmt.Sprintf("wss://ws.coincap.io/prices?assets=%s", strings.Join(assets, ","))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to CoinCap WebSocket for assets: %v\n", assets)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v\n", err)
			break
		}

		fmt.Printf("Message received: %s\n", message)
	}
}

func CoincapInfo() {
	fmt.Println("Coincap API utility initialized")
}
