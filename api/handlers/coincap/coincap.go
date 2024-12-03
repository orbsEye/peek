package coincap

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Coincap struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price_usd,string"`
}

func StartWebSocketServer(assets []string, clientConn *websocket.Conn) {
	url := fmt.Sprintf("wss://ws.coincap.io/prices?assets=%s", strings.Join(assets, ","))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Failed to connect to CoinCap WebSocket: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to CoinCap WebSocket for assets: %v\n", assets)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error from CoinCap WebSocket: %v\n", err)
			break
		}

		var data map[string]string
		if err := json.Unmarshal(message, &data); err != nil {
			log.Printf("JSON unmarshal error: %v\n", err)
			continue
		}

		var result []Coincap
		for symbol, price := range data {
			result = append(result, Coincap{Symbol: symbol, Price: toFloat(price)})
		}

		if err := clientConn.WriteJSON(result); err != nil {
			log.Printf("Error writing JSON to client: %v\n", err)
			break
		}
	}
}

func toFloat(value string) float64 {
	price, _ := strconv.ParseFloat(value, 64)
	return price
}
