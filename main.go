//go:generate protoc --go_out=. proto/card.proto proto/handshake.proto proto/ready.proto proto/changed.proto proto/players.proto proto/turn.proto proto/envelope.proto proto/play.proto
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/caffinatedmonkey/dos/game"
	dosProto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// TODO: More logging
// TODO: Radomness sucks
// TODO: Matchmaking
// TODO: Fix explosive disconnects
var started = utils.NewBroadcaster()
var turnBroadcaster = utils.NewBroadcaster()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var game = dos.NewGame(true)
var commonMessages = utils.NewBroadcaster()

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

	go turnBroadcaster.StartBroadcasting()
	go HandleCommonMessages(game, turnBroadcaster, commonMessages)
	go commonMessages.StartBroadcasting()
	go started.StartBroadcasting()

	fmt.Printf("[server] initializing on %s\n", *listen)
	err := s.ListenAndServe()
	log.Fatal(err)
}

func handleSocket(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("[websocket] new connection")

	rawConn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("[websocket] connection initialization failed", err)
		return
	}

	conn := &LockedSocket{Conn: rawConn}

	handshake := dosProto.HandshakeMessage{}
	err = Read(conn, &handshake)
	if err != nil {
		conn.Close()
		return
	}

	switch handshake.Type {
	case dosProto.ClientType_PLAYER:
		fmt.Println("[websocket] new player client")

		// Wait for player to be ready
		var player *dos.Player

		for {
			ready := dosProto.ReadyMessage{}
			err := ReadMessage(conn, dosProto.MessageType_READY, &ready)
			if err != nil {
				conn.Close()
				return
			}

			player, err = game.NewPlayer(ready.Name)
			if err != nil {
				fmt.Printf("[game] %s failed to join: %v\n", ready.Name, err)
				errorMessage := dosProto.ErrorMessage{Reason: err.Error()}

				err := WriteMessage(conn, dosProto.MessageType_ERROR, &errorMessage)
				if err != nil {
					conn.Close()
					return
				}
			} else {
				fmt.Printf("[game] %s joined\n", ready.Name)
				err := WriteMessage(conn, dosProto.MessageType_SUCCESS, nil)
				if err != nil {
					conn.Close()
					return
				}
				break
			}
		}

		// Handle leaving
		oldHandler := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			// TODO: If its the players turn, handle it

			// Handles a client requested close.
			fmt.Printf("[game] player %s is leaving\n", player.Name)
			game.RemovePlayer(player)

			// Close socket
			conn.WriteLock.Lock()
			ret := oldHandler(code, text)
			conn.WriteLock.Unlock()
			return ret
		})

		// Send player list
		playersMessage := dosProto.PlayersMessage{Additions: game.GetPlayerList()}
		err = WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)
		if err != nil {
			conn.Close()
			return
		}

		messages := make(chan interface{})
		commonMessages.AddListener(messages)

		start := make(chan interface{})
		started.AddListener(start)

		handAdditions := make(chan interface{})
		player.Cards.Additions.AddListener(handAdditions)
		go player.Cards.Additions.StartBroadcasting()

		handDeletions := make(chan interface{})
		player.Cards.Deletions.AddListener(handDeletions)
		go player.Cards.Deletions.StartBroadcasting()

		turnChanged := make(chan interface{})
		turnBroadcaster.AddListener(turnChanged)

		started := false
		isMyTurn := false

		hasDrawn := false
		hasPlayed := false

		go func() {
			for {
				var buf []byte
				var err error

				select {
				case message := <-messages:
					var ok bool
					buf, ok = message.([]byte)
					if !ok {
						log.Println("[error] commonMessages is broadcasting a non []byte value")
						continue
					}

				case <-start:
					additions := make([]*dosProto.Card, len(player.Cards.List))
					for index := range player.Cards.List {
						additions[index] = &player.Cards.List[index]
					}
					changed := &dosProto.CardsChangedMessage{Additions: additions}
					buf, err = ZipMessage(dosProto.MessageType_CARDS, changed)
					started = true

				case deletion := <-handDeletions:
					if !started {
						continue
					}

					msg := &dosProto.CardsChangedMessage{
						Deletions: []int32{deletion.(int32)},
					}
					buf, err = ZipMessage(dosProto.MessageType_CARDS, msg)

				case addition := <-handAdditions:
					if !started {
						continue
					}

					card := addition.(dosProto.Card)
					msg := &dosProto.CardsChangedMessage{
						Additions: []*dosProto.Card{&card},
					}
					buf, err = ZipMessage(dosProto.MessageType_CARDS, msg)

				case turn := <-turnChanged:
					isMyTurn = player.Name == turn.(string)
					hasDrawn = false
					hasPlayed = false
					continue
				}

				if err != nil {
					log.Println("[protobuf] encoding error", err)
					continue
				}

				err = conn.WriteMessage(websocket.BinaryMessage, buf)
				if err != nil {
					log.Println("[websocket] write error", err)
				}
			}
		}()

		for {
			envelope := dosProto.Envelope{}
			err := Read(conn, &envelope)
			if err != nil {
				conn.Close()
				return
			}

			if !isMyTurn {
				// Ignore messages sent during other people's turns
				return
			}

			switch envelope.Type {
			case dosProto.MessageType_DRAW:
				if !hasDrawn && !hasPlayed {
					fmt.Printf("[game] %s drawing card\n", player.Name)
					game.DrawCards(&player.Cards, 1)
					hasDrawn = true
				}

			case dosProto.MessageType_PLAY:
				if !hasPlayed {
					fmt.Printf("[game] %s playing card\n", player.Name)

					playMessage := dosProto.PlayMessage{}
					err := proto.Unmarshal(envelope.Contents, &playMessage)
					if err != nil {
						fmt.Println("[protobuf] failed to parse message:", err)
						return
					}

					err = game.PlayCard(player, playMessage.Id, playMessage.Color)
					if err != nil {
						fmt.Println("[game] play failed:", err, playMessage)
					} else {
						fmt.Println("[game] played card")
						hasPlayed = true
					}
				}

			case dosProto.MessageType_DONE:
				fmt.Println("done message sent")
				if hasDrawn || hasPlayed {
					fmt.Printf("[game] %s is done with turn\n", player.Name)
					player.TurnDone <- struct{}{}
				}
			}
		}

	case dosProto.ClientType_SPECTATOR:
		fmt.Println("[game] spectator joined")

		// Send player list
		playersMessage := dosProto.PlayersMessage{
			Additions: game.GetPlayerList(),
		}
		err = WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)
		if err != nil {
			conn.Close()
			return
		}

		messages := make(chan interface{})
		commonMessages.AddListener(messages)

		go func() {
			for {
				var buf []byte
				var err error

				select {
				case message := <-messages:
					var ok bool
					buf, ok = message.([]byte)
					if !ok {
						log.Println("[error] commonMessages is broadcasting a non []byte value")
						continue
					}
				}

				err = conn.WriteMessage(websocket.BinaryMessage, buf)
				if err != nil {
					log.Println("[websocket] write error", err)
				}
			}
		}()

		envelope := dosProto.Envelope{}
		err := Read(conn, &envelope)
		if err != nil {
			conn.Close()
			return
		}

		if envelope.Type != dosProto.MessageType_START {
			// Start is the only message spectators can send.
			return
		}

		fmt.Println("[game] spectator starting game")
		started.Broadcast(nil)

		for {
			player := game.NextPlayer()
			fmt.Printf("[game] player %s's turn\n", player.Name)
			turnBroadcaster.Broadcast(player.Name)
			<-player.TurnDone

			if len(player.Cards.List) == 0 {
				// Game Done!
				// TODO:
				//  Send message to players and spectators
				//  Cleanup
				return
			}
		}
	}
}

func HandleCommonMessages(game *dos.Game, turnBroadcaster *utils.Broadcaster, outputMessages *utils.Broadcaster) {
	turnChannel := make(chan interface{})
	turnBroadcaster.AddListener(turnChannel)

	for {
		var err error
		var bytes []byte

		select {
		case newPlayer := <-game.PlayerJoined:
			msg := &dosProto.PlayersMessage{
				Additions: []string{newPlayer},
			}
			bytes, err = ZipMessage(dosProto.MessageType_PLAYERS, msg)

		case leavingPlayer := <-game.PlayerLeft:
			msg := &dosProto.PlayersMessage{
				Deletions: []string{leavingPlayer},
			}
			bytes, err = ZipMessage(dosProto.MessageType_PLAYERS, msg)

		case nextPlayer := <-turnChannel:
			lastCard := game.Discard.List[len(game.Discard.List)-1]
			msg := &dosProto.TurnMessage{
				LastPlayed: &lastCard,
				Player:     nextPlayer.(string),
			}
			bytes, err = ZipMessage(dosProto.MessageType_TURN, msg)
		}

		if err != nil {
			// Encoding Error
			log.Println("[protobuf] Encoding error", err)
		} else {
			outputMessages.Broadcast(bytes)
		}
	}
}

func Read(conn *LockedSocket, message proto.Message) error {
	conn.ReadLock.Lock()
	format, buf, err := conn.ReadMessage()
	conn.ReadLock.Unlock()
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

func ReadMessage(conn *LockedSocket, typ dosProto.MessageType, message proto.Message) error {
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

func Write(conn *LockedSocket, message proto.Message) error {
	buf, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("[websocket] failed to encode message:", err)
		return err
	}

	conn.WriteLock.Lock()
	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	conn.WriteLock.Unlock()
	if err != nil {
		fmt.Println("[websocket] failed to write message:", err)
		return err
	}

	return nil
}

func WriteMessage(conn *LockedSocket, typ dosProto.MessageType, message proto.Message) error {
	buf, err := ZipMessage(typ, message)
	if err != nil {
		fmt.Println("[composing] failed to compose message:", err)
		return err
	}

	conn.WriteLock.Lock()
	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	conn.WriteLock.Unlock()

	if err != nil {
		fmt.Println("[websocket] failed to write message:", err)
		return err
	}

	return nil
}

func ZipMessage(typ dosProto.MessageType, message proto.Message) ([]byte, error) {
	var buf []byte
	var err error
	if message != nil {
		buf, err = proto.Marshal(message)
		if err != nil {
			return nil, err
		}
	}

	envelope := dosProto.Envelope{}
	envelope.Type = typ
	envelope.Contents = buf

	return proto.Marshal(&envelope)
}

type LockedSocket struct {
	ReadLock  sync.Mutex
	WriteLock sync.Mutex
	*websocket.Conn
}
