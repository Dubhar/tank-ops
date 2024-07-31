package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	var gameState *GameState

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if msgType != websocket.TextMessage {
			log.Println("Wrong message type.")
			break
		}

		log.Printf("player messaage %s\n", msg)

		clientMsg := ClientMessage{}
		json.Unmarshal(msg, &clientMsg)

		time.Sleep(time.Millisecond * 200)

		switch clientMsg.Type {
		case ClientStartGame:
			gameState = newGameState(basicConfig)
			err := conn.WriteJSON(newStartGameMessage(basicConfig))
			if err != nil {
				panic(err)
			}
			continue
		case ClientSendTurn:
			if gameState == nil {
				panic("gs is nil")
			}
			log.Println(clientMsg.Actions)

			results := gameState.resolveActions(clientMsg.Actions)
			// time.Sleep(time.Millisecond * 1300)
			time.Sleep(time.Millisecond * 150)
			err := conn.WriteJSON(newTurnResultsMessage(results))
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	fmt.Println("hello")

	static := http.Dir("web/dist")

	http.Handle("/", http.FileServer(static))
	http.HandleFunc("/ws", serveWS)

	http.ListenAndServe("localhost:8000", nil)
}
