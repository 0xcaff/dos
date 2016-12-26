//go:generate protoc --go_out=. proto/card.proto proto/handshake.proto proto/ready.proto proto/changed.proto proto/players.proto proto/turn.proto proto/envelope.proto proto/play.proto proto/session.proto
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/caffinatedmonkey/dos/game"
	dosProto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var GameStore = make(map[int32]*GameState)
var StoreMutex sync.Mutex

var (
	listen   = flag.String("listen", ":8080", "Address to serve on")
	upgrader = websocket.Upgrader{
		// TODO: Implement
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	flag.Parse()

	// TODO: Doesn't seed for every random
	rand.Seed(time.Now().Unix())

	fs := SinglePageFileSystem{rice.MustFindBox("frontend/build").HTTPBox()}

	mux := http.NewServeMux()
	mux.HandleFunc("/socket", handleSocket)
	mux.Handle("/", http.FileServer(fs))

	s := &http.Server{
		Addr:    *listen,
		Handler: mux,
	}

	log.Printf("[server] initializing on %s\n", *listen)
	err := s.ListenAndServe()
	log.Fatal(err)
}

func handleSocket(rw http.ResponseWriter, r *http.Request) {
	log.Println("[websocket] new connection")

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("[websocket] connection initialization failed", err)
		return
	}

	handshake := dosProto.HandshakeMessage{}
	err = Read(conn, &handshake)
	if err != nil {
		return
	}

	switch handshake.Type {
	case dosProto.ClientType_PLAYER:
		handlePlayer(conn)

	case dosProto.ClientType_SPECTATOR:
		handleSpectator(conn)
	}
}

type GameState struct {
	Game           *dos.Game
	Started        utils.Broadcaster
	Turn           utils.Broadcaster
	CommonMessages utils.Broadcaster
	IsStarted      bool
}

// Encodes messages common to all clients and broadcasts them on CommonMessages.
func (state *GameState) HandleMessages() {
	turnChannel := state.Turn.NewListener()

	for {
		var err error
		var bytes []byte

		select {
		case newPlayer := <-state.Game.PlayerJoined:
			msg := &dosProto.PlayersMessage{
				Additions: []string{newPlayer},
			}
			bytes, err = ZipMessage(dosProto.MessageType_PLAYERS, msg)

		case leavingPlayer := <-state.Game.PlayerLeft:
			msg := &dosProto.PlayersMessage{
				Deletions: []string{leavingPlayer},
			}
			bytes, err = ZipMessage(dosProto.MessageType_PLAYERS, msg)

		case nextPlayer := <-turnChannel:
			lastCard := state.Game.Discard.List[len(state.Game.Discard.List)-1]
			msg := &dosProto.TurnMessage{
				LastPlayed: &lastCard,
				Player:     nextPlayer.(string),
			}
			bytes, err = ZipMessage(dosProto.MessageType_TURN, msg)
		}

		if err != nil {
			log.Println("[protobuf] Encoding error", err)
		} else {
			state.CommonMessages.Broadcast(bytes)
		}
	}
}

// TODO: Better logging with logrus
func handleSpectator(conn *websocket.Conn) {
	log.Println("[websocket] spectator joined")

	// Get an unused session id
	var session int32
	for ok := true; ok; _, ok = GameStore[session] {
		session = rand.Int31n(1000000)
	}

	game := dos.NewGame(true)
	state := &GameState{
		Game:           game,
		Started:        *utils.NewBroadcaster(),
		Turn:           *utils.NewBroadcaster(),
		CommonMessages: *utils.NewBroadcaster(),
		IsStarted:      false,
	}

	go state.Turn.StartBroadcasting()
	go state.CommonMessages.StartBroadcasting()
	go state.Started.StartBroadcasting()
	go state.HandleMessages()

	StoreMutex.Lock()
	GameStore[session] = state
	StoreMutex.Unlock()

	defer func() {
		// TODO: Disconnect Players
		StoreMutex.Lock()
		delete(GameStore, session)
		StoreMutex.Unlock()
	}()

	go func() {
		// Send Session ID For Display
		sessionMessage := dosProto.SessionMessage{Session: session}
		err := WriteMessage(conn, dosProto.MessageType_SESSION, &sessionMessage)
		if err != nil {
			return
		}

		messages := state.CommonMessages.NewListener()
		defer state.CommonMessages.RemoveListener(messages)

		for {
			var buf []byte
			var err error

			message := <-messages
			buf = message.([]byte)

			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.Println("[websocket] failed to write", err)
				return
			}
		}
	}()

	// Start is the only message spectators can send.
	err := ReadMessage(conn, dosProto.MessageType_START, nil)
	if err != nil {
		return
	}

	log.Println("[game] starting")

	state.IsStarted = true
	state.Started.Broadcast(nil)

	for {
		player := game.NextPlayer()
		log.Printf("[game] (%s) turn\n", player.Name)
		state.Turn.Broadcast(player.Name)

		select {
		case <-player.TurnDone:
			log.Printf("[game] (%s) turn done\n", player.Name)
			if len(player.Cards.List) == 0 {
				log.Printf("[game] (%s) done with game\n", player.Name)

				// TODO: Send to players.
				// conn.WriteControl(
				// 	websocket.CloseMessage,
				// 	websocket.FormatCloseMessage(websocket.CloseNormalClosure, "won"),
				// 	time.Now().Add(time.Second),
				// )

				game.RemovePlayer(player)
			}
		}
	}
}

func handlePlayer(conn *websocket.Conn) {
	log.Println("[websocket] new player joined")

	// TODO: limit this and kill connection when it goes over.
	var state *GameState
	for {
		sessionMessage := &dosProto.SessionMessage{}
		err := ReadMessage(conn, dosProto.MessageType_SESSION, sessionMessage)
		if err != nil {
			return
		}

		var ok bool
		state, ok = GameStore[sessionMessage.Session]
		if !ok {
			errorMessage := dosProto.ErrorMessage{
				Reason: "Invalid Game PIN. That game doesn't exist.",
			}

			err := WriteMessage(conn, dosProto.MessageType_ERROR, &errorMessage)
			if err != nil {
				return
			}
		} else {
			err := WriteMessage(conn, dosProto.MessageType_SUCCESS, nil)
			if err != nil {
				return
			}
			break
		}
	}
	game := state.Game

	// Wait for player to be ready
	var player *dos.Player

	// TODO: limit this
	for {
		ready := dosProto.ReadyMessage{}
		err := ReadMessage(conn, dosProto.MessageType_READY, &ready)
		if err != nil {
			return
		}

		player, err = game.NewPlayer(ready.Name)
		if err != nil {
			log.Printf("[game] (%s) failed to join: %v\n", ready.Name, err)
			errorMessage := dosProto.ErrorMessage{Reason: err.Error()}

			err := WriteMessage(conn, dosProto.MessageType_ERROR, &errorMessage)
			if err != nil {
				return
			}

		} else {
			log.Printf("[game] (%s) joined\n", ready.Name)
			err := WriteMessage(conn, dosProto.MessageType_SUCCESS, nil)
			if err != nil {
				return
			}
			break
		}
	}

	// Handle leaving
	oldHandler := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("[game] (%s) is leaving\n", player.Name)
		// Close socket
		ret := oldHandler(code, text)

		game.RemovePlayer(player)
		player.TurnDone <- struct{}{}

		return ret
	})

	// Send player list
	playersMessage := dosProto.PlayersMessage{Additions: game.GetPlayerList()}
	err := WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)
	if err != nil {
		return
	}

	isStarted := false
	isMyTurn := false

	hasDrawn := false
	hasPlayed := false

	go func() {
		handAdditions := player.Cards.Additions.NewListener()
		go player.Cards.Additions.StartBroadcasting()
		defer player.Cards.Additions.Destroy()

		handDeletions := player.Cards.Deletions.NewListener()
		go player.Cards.Deletions.StartBroadcasting()
		defer player.Cards.Additions.Destroy()

		messages := state.CommonMessages.NewListener()
		defer state.CommonMessages.RemoveListener(messages)

		start := state.Started.NewListener()
		defer state.Started.RemoveListener(start)

		turnChanged := state.Turn.NewListener()
		defer state.Turn.RemoveListener(turnChanged)

		for {
			var buf []byte
			var err error

			select {
			case message := <-messages:
				buf = message.([]byte)

			case <-start:
				additions := make([]*dosProto.Card, len(player.Cards.List))
				for index := range player.Cards.List {
					additions[index] = &player.Cards.List[index]
				}
				changed := &dosProto.CardsChangedMessage{Additions: additions}
				buf, err = ZipMessage(dosProto.MessageType_CARDS, changed)
				isStarted = true

			case deletion := <-handDeletions:
				if !isStarted {
					continue
				}

				msg := &dosProto.CardsChangedMessage{
					Deletions: []int32{deletion.(int32)},
				}
				buf, err = ZipMessage(dosProto.MessageType_CARDS, msg)

			case addition := <-handAdditions:
				if !isStarted {
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

				// TODO: GameDone Broadcaster
			}

			if err != nil {
				log.Println("[protobuf] encoding error", err)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.Println("[websocket] failed to write message", err)
				return
			}
		}
	}()

	for {
		envelope := dosProto.Envelope{}
		err := Read(conn, &envelope)
		if err != nil {
			return
		}

		if !isMyTurn {
			// Ignore messages sent during other people's turns
			return
		}

		switch envelope.Type {
		case dosProto.MessageType_DRAW:
			if !hasDrawn && !hasPlayed {
				log.Printf("[game] (%s) drawing card\n", player.Name)
				game.DrawCards(&player.Cards, 1)
				hasDrawn = true
			}

		case dosProto.MessageType_PLAY:
			if !hasPlayed {
				log.Printf("[game] (%s) playing card\n", player.Name)

				playMessage := dosProto.PlayMessage{}
				err := proto.Unmarshal(envelope.Contents, &playMessage)
				if err != nil {
					log.Println("[protobuf] failed to parse message:", err)
					conn.WriteControl(
						websocket.CloseMessage,
						websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
						time.Now().Add(time.Second),
					)
					return
				}

				err = game.PlayCard(player, playMessage.Id, playMessage.Color)
				if err != nil {
					log.Printf("[game] (%s) tried playing card and failed: %#v %#v\n", player.Name, err, playMessage)
				} else {
					log.Printf("[game] (%s) played card\n", player.Name)
					hasPlayed = true
				}
			}

		case dosProto.MessageType_DONE:
			if hasDrawn || hasPlayed {
				log.Printf("[game] (%s) done with turn\n", player.Name)
				player.TurnDone <- struct{}{}
			}
		}
	}
}

func Read(conn *websocket.Conn, message proto.Message) error {
	format, buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if format != websocket.BinaryMessage {
		log.Println("[websocket] got non binary message from", conn.RemoteAddr())
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)

		return fmt.Errorf("dos: got non binary message\n")
	}

	err = proto.Unmarshal(buf, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)

		return fmt.Errorf("[protobuf] failed to parse message: %#v", err)
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
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)
		err = fmt.Errorf("Received type %s instead of type %s", envelope.Type.String(), typ.String())
		log.Println("[websocket]", err)
		return err
	}

	if message == nil {
		return nil
	}

	err = proto.Unmarshal(envelope.Contents, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)

		err = fmt.Errorf("[protobuf] failed to parse envelope: %#v", err)
		log.Println("[websocket]", err)
		return err
	}

	return nil
}

func WriteMessage(conn *websocket.Conn, typ dosProto.MessageType, message proto.Message) error {
	buf, err := ZipMessage(typ, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""),
			time.Now().Add(time.Second),
		)

		log.Println("[protobuf] failed to compose message:", err)
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		log.Println("[websocket] failed to write message:", err)
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
