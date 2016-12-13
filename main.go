package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caffinatedmonkey/dos/game"
	dosProto "github.com/caffinatedmonkey/dos/proto"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var game = dos.NewGame()

var (
	listen = flag.String("listen", ":8080", "Address to serve on")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/socket", handleSocket)
	mux.Handle("/", http.FileServer(SPAFileSystem("frontend/build")))

	s := &http.Server{
		Addr:    *listen,
		Handler: mux,
	}

	errChan := make(chan error)
	go func() {
		fmt.Printf("[server] initializing on %s\n", *listen)
		err := s.ListenAndServe()
		errChan <- err
	}()

	err := <-errChan
	log.Fatal(err)
}

func handleSocket(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("[websocket] new connection")

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("[websocket] connection initialization failed", err)
		return
	}

	handshake := dosProto.HandshakeMessage{}
	Read(conn, &handshake)

	switch handshake.Type {
	case dosProto.ClientType_PLAYER:
		fmt.Println("[websocket] new player client")

		// Wait for player to be ready
		var name string
		for {
			ready := dosProto.ReadyMessage{}
			err := ReadMessage(conn, dosProto.MessageType_READY, &ready)
			if err != nil {
				fmt.Println("[websocket] failed to parse message", err)
				return
			}

			_, err = game.NewPlayer(ready.Name)
			if err != nil {
				fmt.Printf("[game] %s failed to join: %v\n", ready.Name, err)
				// TODO: Send error downstream
			} else {
				fmt.Printf("[game] %s joined\n", ready.Name)
				name = ready.Name
				break
			}
		}

		// Handle leaving
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Printf("[game] player %s is leaving\n", name)
			game.RemovePlayer(name)
			return nil
		})

		// Send player list
		playersMessage := dosProto.PlayersMessage{}
		playersMessage.Initial = game.GetPlayerList()
		WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)

		// Maintain player list
		go SendPlayerJoins(conn, game)
		go SendPlayerLeaves(conn, game)

		conn.ReadMessage()

	case dosProto.ClientType_SPECTATOR:
		fmt.Println("[game] spectator joined")

		// Send player list
		playersMessage := dosProto.PlayersMessage{}
		playersMessage.Initial = game.GetPlayerList()
		WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)

		// Maintain player list
		go SendPlayerJoins(conn, game)
		go SendPlayerLeaves(conn, game)

		envelope := dosProto.Envelope{}
		Read(conn, &envelope)

		if envelope.Type != dosProto.MessageType_START {
			// Start is the only message spectators can send.
			return
		}

		fmt.Println("[game] spectator starting game")

		// TODO: Handle start
		// Tell players about their cards
		// Pick a random player to use as the first one
		// Send a turn message
	}
}

func SendPlayerJoins(conn *websocket.Conn, game *dos.Game) {
	joined := make(chan interface{})
	game.PlayerJoined.AddListener(joined)
	go game.PlayerJoined.StartBroadcasting()

	for newPlayerName := range joined {
		msg := dosProto.PlayersMessage{}
		msg.Addition = newPlayerName.(string)

		err := WriteMessage(conn, dosProto.MessageType_PLAYERS, &msg)
		if err != nil {
			// Socket is closed/something bad happened
			game.PlayerJoined.RemoveListener(joined)
			close(joined)
			conn.Close()
		}
	}
}

func SendPlayerLeaves(conn *websocket.Conn, game *dos.Game) {
	left := make(chan interface{})
	game.PlayerLeft.AddListener(left)
	go game.PlayerLeft.StartBroadcasting()

	for leavingPlayer := range left {
		msg := dosProto.PlayersMessage{}
		msg.Deletion = leavingPlayer.(string)

		err := WriteMessage(conn, dosProto.MessageType_PLAYERS, &msg)
		if err != nil {
			// Socket is closed/something bad happened
			game.PlayerJoined.RemoveListener(left)
			close(left)
			conn.Close()
		}
	}
}

func Read(conn *websocket.Conn, message proto.Message) error {
	format, buf, err := conn.ReadMessage()
	if format != websocket.BinaryMessage {
		fmt.Println("[websocket] warning! reading non binary message")
	}

	if err != nil {
		fmt.Println("[websocket] failed to read message:", err)
		return err
	}

	err = proto.Unmarshal(buf, message)
	if err != nil {
		fmt.Println("[websocket] failed to parse handshake:", err)
		return err
	}

	return nil
}

func ReadMessage(conn *websocket.Conn, typ dosProto.MessageType, message proto.Message) error {
	envelope := dosProto.Envelope{}
	err := Read(conn, &envelope)
	if err != nil {
		return err
	}

	if envelope.Type != typ {
		return fmt.Errorf("Received type %s instead of type %s", envelope.Type.String(), typ.String())
	}

	return proto.Unmarshal(envelope.Contents, message)
	if err != nil {
		return err
	}

	return nil
}

func Write(conn *websocket.Conn, message proto.Message) error {
	buf, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("[websocket] failed to encode message:", err)
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		fmt.Println("[websocket] failed to write message:", err)
		return err
	}

	return nil
}

func WriteMessage(conn *websocket.Conn, typ dosProto.MessageType, message proto.Message) error {
	buf, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("[websocket] failed to encode message", err)
		return err
	}

	envelope := dosProto.Envelope{}
	envelope.Type = typ
	envelope.Contents = buf

	return Write(conn, &envelope)
}
