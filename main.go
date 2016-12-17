package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caffinatedmonkey/dos/game"
	dosProto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// TODO: Websocket concurrent writes aren't allowed

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
			return oldHandler(code, text)
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
		Read(conn, &envelope)

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

			// TODO: Check if game is done
		}
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
			game.PlayerLeft.RemoveListener(left)
			close(left)
			conn.Close()
		}
	}
}

func SendCardAdditions(conn *websocket.Conn, cards *dos.Cards) {
	additions := make(chan interface{})
	cards.Additions.AddListener(additions)
	go cards.Additions.StartBroadcasting()

	for addition := range additions {
		card := addition.(dosProto.Card)
		msg := dosProto.CardsChangedMessage{}
		msg.Additions = []*dosProto.Card{&card}
		WriteMessage(conn, dosProto.MessageType_CARDS, &msg)
	}
}

func SendCardDeletions(conn *websocket.Conn, cards *dos.Cards) {
	deletions := make(chan interface{})
	cards.Deletions.AddListener(deletions)
	go cards.Deletions.StartBroadcasting()

	for deletion := range deletions {
		msg := dosProto.CardsChangedMessage{}
		msg.Deletions = []int32{deletion.(int32)}
		WriteMessage(conn, dosProto.MessageType_CARDS, &msg)
	}
}

func SendTurnChanged(conn *websocket.Conn, game *dos.Game) {
	nextTurn := make(chan interface{})
	game.Turn.AddListener(nextTurn)
	go game.Turn.StartBroadcasting()

	for turn := range nextTurn {
		lastCard := game.Discard.List[len(game.Discard.List)-1]
		msg := dosProto.TurnMessage{}
		msg.LastPlayed = &lastCard
		msg.Player = turn.(string)
		WriteMessage(conn, dosProto.MessageType_TURN, &msg)
	}
}

func HandleTurn(conn *websocket.Conn, player *dos.Player, game *dos.Game) {
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
				Read(conn, &envelope)

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
