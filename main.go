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

// TODO: Radomness sucks
// TODO: Matchmaking
// TODO: When a player plays a card before their turn, when it is their turn the
// card is played. This shouldn't happen.
// TODO: Fix explosive disconnects
var started = utils.NewBroadcaster()

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

	go started.StartBroadcasting()

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

	rawConn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("[websocket] connection initialization failed", err)
		return
	}

	conn := &LockedSocket{Conn: rawConn}

	handshake := dosProto.HandshakeMessage{}
	err := Read(conn, &handshake)
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
				fmt.Println("[websocket] failed to parse message", err)
				return
			}

			player, err = game.NewPlayer(ready.Name)
			if err != nil {
				fmt.Printf("[game] %s failed to join: %v\n", ready.Name, err)
				// TODO: Send error downstream
			} else {
				fmt.Printf("[game] %s joined\n", ready.Name)
				break
			}
		}

		// Handle leaving
		oldHandler := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			// TODO: Not triggered on closes without an active read
			fmt.Printf("[game] player %s is leaving\n", player.Name)
			game.RemovePlayer(player)

			// Close socket
			conn.WriteLock.Lock()
			ret := oldHandler(code, text)
			conn.WriteLock.Unlock()
			return ret
		})

		// Send player list
		playersMessage := dosProto.PlayersMessage{}
		playersMessage.Initial = game.GetPlayerList()
		WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)

		// Maintain player list
		go SendPlayerJoins(conn, game)
		go SendPlayerLeaves(conn, game)

		// Handle turn
		go SendTurnChanged(conn, game)
		go HandleTurn(conn, player, game)

		// Wait for game start
		start := make(chan interface{})
		started.AddListener(start)
		<-start

		// Synchronize cards
		changed := dosProto.CardsChangedMessage{}
		additions := make([]*dosProto.Card, len(player.Cards.List))
		for index := range player.Cards.List {
			additions[index] = &player.Cards.List[index]
		}
		changed.Additions = additions
		WriteMessage(conn, dosProto.MessageType_CARDS, &changed)

		go SendCardAdditions(conn, &player.Cards)
		go SendCardDeletions(conn, &player.Cards)

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
		go SendTurnChanged(conn, game)
		started.Broadcast(nil)

		for {
			player := game.NextPlayer()
			fmt.Printf("[game] player %s's turn\n", player.Name)
			game.Turn.Broadcast(player.Name)
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

func SendPlayerJoins(conn *LockedSocket, game *dos.Game) {
	WriteMessageOn(conn, &game.PlayerJoined,
		func(conn *LockedSocket, newPlayerName interface{}) (proto.Message, error) {
			msg := dosProto.PlayersMessage{}
			msg.Addition = newPlayerName.(string)
			return ZipMessage(dosProto.MessageType_PLAYERS, &msg)
		})
}

func SendPlayerLeaves(conn *LockedSocket, game *dos.Game) {
	WriteMessageOn(conn, &game.PlayerLeft,
		func(conn *LockedSocket, leavingPlayer interface{}) (proto.Message, error) {
			msg := dosProto.PlayersMessage{}
			msg.Deletion = leavingPlayer.(string)
			return ZipMessage(dosProto.MessageType_PLAYERS, &msg)
		})
}

func SendCardAdditions(conn *LockedSocket, cards *dos.Cards) {
	WriteMessageOn(conn, &cards.Additions,
		func(conn *LockedSocket, addition interface{}) (proto.Message, error) {
			card := addition.(dosProto.Card)
			msg := dosProto.CardsChangedMessage{}
			msg.Additions = []*dosProto.Card{&card}
			return ZipMessage(dosProto.MessageType_CARDS, &msg)
		})
}

func SendCardDeletions(conn *LockedSocket, cards *dos.Cards) {
	WriteMessageOn(conn, &cards.Deletions,
		func(conn *LockedSocket, deletion interface{}) (proto.Message, error) {
			msg := dosProto.CardsChangedMessage{}
			msg.Deletions = []int32{deletion.(int32)}
			return ZipMessage(dosProto.MessageType_CARDS, &msg)
		})
}

func SendTurnChanged(conn *LockedSocket, game *dos.Game) {
	WriteMessageOn(conn, &game.Turn,
		func(conn *LockedSocket, turn interface{}) (proto.Message, error) {
			lastCard := game.Discard.List[len(game.Discard.List)-1]
			msg := dosProto.TurnMessage{}
			msg.LastPlayed = &lastCard
			msg.Player = turn.(string)
			return ZipMessage(dosProto.MessageType_TURN, &msg)
		})
}

func HandleTurn(conn *LockedSocket, player *dos.Player, game *dos.Game) {
	nextTurn := make(chan interface{})
	game.Turn.AddListener(nextTurn)
	go game.Turn.StartBroadcasting()

	// TODO: This blocks the broadcaster so it should be run last or the
	// broadcaster should be buffered.

	for turn := range nextTurn {
		if turn.(string) == player.Name {
			lastCard := game.Discard.List[len(game.Discard.List)-1]
			fmt.Printf("[game] top of discard pile %#v\n", lastCard)

			done := false
			hasDrawn := false
			hasPlayed := false

			for !done {
				envelope := dosProto.Envelope{}
				err := Read(conn, &envelope)
				if err != nil {
					conn.Close()
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
							// TODO: Handle error
						} else {
							fmt.Println("[game] played card")
							hasPlayed = true
						}
					}

				case dosProto.MessageType_DONE:
					if hasDrawn || hasPlayed {
						fmt.Printf("[game] %s is done with turn\n", player.Name)
						done = true
						player.TurnDone <- nil
					}
				}
			}
		}
	}
}

func WriteMessageOn(conn *LockedSocket, broadcaster *utils.Broadcaster,
	composer func(conn *LockedSocket, thing interface{}) (proto.Message, error)) {

	channel := make(chan interface{})
	broadcaster.AddListener(channel)
	go broadcaster.StartBroadcasting()

	for thing := range channel {
		message, err := composer(conn, thing)
		if err != nil {
			fmt.Println("[composer] error while composing message:", err)
			continue
		}

		err = Write(conn, message)
		if err != nil {
			// Socket is closed/something bad happened
			broadcaster.RemoveListener(channel)
			close(channel)
			conn.Close()
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
	envelope, err := ZipMessage(typ, message)
	if err != nil {
		fmt.Println("[composing] failed to compose message:", err)
		return err
	}

	return Write(conn, envelope)
}

func ZipMessage(typ dosProto.MessageType, message proto.Message) (proto.Message, error) {
	buf, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	envelope := dosProto.Envelope{}
	envelope.Type = typ
	envelope.Contents = buf

	return &envelope, nil
}

type LockedSocket struct {
	ReadLock  sync.Mutex
	WriteLock sync.Mutex
	*websocket.Conn
}
