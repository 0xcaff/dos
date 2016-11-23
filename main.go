package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/caffinatedmonkey/dos/game"
	"github.com/caffinatedmonkey/dos/proto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var wg = sync.WaitGroup{}
var mux = http.NewServeMux()
var upgrader = websocket.Upgrader{}
var s = &http.Server{
	Addr:    ":8080",
	Handler: mux,
}
var game = dos.NewGame()

func main() {
	mux.HandleFunc("/ws", handleSocket)
	mux.Handle("/", http.FileServer(SPAFileSystem("frontend")))

	wg.Add(1)
	go func(s *http.Server) {
		fmt.Println("[server] initializing")
		log.Fatal(s.ListenAndServe())
		wg.Done()
	}(s)

	wg.Wait()
}

func handleSocket(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("[websocket] connection initialization failed", err)
		return
	}

	_, buf, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("[websocket] failed to read message", err)
		return
	}

	handshake := dos.Handshake{}
	err = proto.Unmarshal(handshake, buf)
	if err != nil {
		fmt.Println("[websocket] failed to parse handshake", err)
		return
	}

	// TODO: Behave based on handshake.
	switch handshake.Type {
	case HandshakeMessage_PLAYER:

	case HandshakeMessage_DECK:
	}
}

func handlePlay(conn *websocket.Conn) {
	_, buf, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("[websocket] failed to read message", err)
		return
	}

	readyMessage := dos.ReadyMessage{}
	err = proto.Unmarshal(readyMessage, buf)
	if err != nil {
		fmt.Println("[websocket] failed to ready message", err)
		return
	}

	player := game.NewPlayer(readyMessage.Name)
	fmt.Printf("[game] %s joined\n", readyMessage.Name)

	// TODO: Setup listeners to convey state
}

func handleDeck(conn *websocket.Conn) {

}

func handleSpectate(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("[websocket] connection failed")
		return
	}

	fmt.Println("[game] spectator joined")
	game.NewSpectator(*conn)
	err = conn.WriteJSON(game)
	if err != nil {
		fmt.Printf("[websocket] %v\n", err)
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[websocket] %v\n", err)
			break
		}

		if messageType == websocket.TextMessage && string(p) == "start" {
			fmt.Println("[game] starting game")
			game.Start()
		}
	}
}
