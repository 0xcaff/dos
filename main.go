//go:generate protoc --go_out=. proto/card.proto proto/handshake.proto proto/ready.proto proto/changed.proto proto/players.proto proto/turn.proto proto/envelope.proto proto/play.proto proto/session.proto
package main

import (
	"flag"
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

// Maximum failed attempts to join a session or select a name before the
// connection is terminated.
const MaxAttempts = 10

var GameStore = make(map[int32]*GameState)
var StoreMutex sync.Mutex

var (
	listen   = flag.String("listen", ":8080", "Address to serve on")
	upgrader = websocket.Upgrader{
		// The reverse proxy will have to enfore origin policies.
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	flag.Parse()

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
	log.Println(err)
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
	Game *dos.Game

	// Encoded []byte sent to all clients (players and spectators) connected to
	// the game.
	CommonMessages utils.Broadcaster

	// Broadcasted by the spectator when the game is started.
	Started   utils.Broadcaster
	IsStarted bool

	// Name (as a string) of the next player broadcasted by the spectator.
	Turn utils.Broadcaster

	// Notify players when they should disconnect.
	PlayerDone utils.Broadcaster

	// A message is sent on this broadcaster when the spectator leaves the game.
	// All players should disconnect when this occurs.
	SpectatorDone utils.Broadcaster
}

func NewGame() (*GameState, int32) {
	// Get an unused session id
	var session int32
	for ok := true; ok; _, ok = GameStore[session] {
		session = rand.Int31n(1000000)
	}

	state := &GameState{
		Game:           dos.NewGame(true),
		CommonMessages: *utils.NewBroadcaster(),
		Started:        *utils.NewBroadcaster(),
		IsStarted:      false,
		Turn:           *utils.NewBroadcaster(),
		PlayerDone:     *utils.NewBroadcaster(),
		SpectatorDone:  *utils.NewBroadcaster(),
	}

	go state.Started.StartBroadcasting()
	go state.Turn.StartBroadcasting()

	go state.PlayerDone.StartBroadcasting()
	go state.SpectatorDone.StartBroadcasting()

	go state.CommonMessages.StartBroadcasting()
	go state.HandleMessages()

	StoreMutex.Lock()
	GameStore[session] = state
	StoreMutex.Unlock()

	return state, session
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

// Removes the state from the global store.
func (state *GameState) Destroy(id int32) {
	StoreMutex.Lock()
	delete(GameStore, id)
	StoreMutex.Unlock()
}

type CloseMessage struct {
	// The player for which this close message is for. The message is handled
	// by this player.
	Name string

	// Code passed to websocket.FormatCloseMessage
	Code int

	// Text passed to websocket.FormatCloseMessage.
	Text string
}

// TODO: Better logging with logrus
func handleSpectator(conn *websocket.Conn) {
	log.Println("[websocket] spectator joined")

	state, session := NewGame()
	game := state.Game
	go func() {
		// Send Session ID For Display
		sessionMessage := dosProto.SessionMessage{Session: session}
		err := WriteMessage(conn, dosProto.MessageType_SESSION, &sessionMessage)
		if err != nil {
			return
		}

		// Forward Common Messages
		messages := state.CommonMessages.NewListener()
		defer state.CommonMessages.RemoveListener(messages)

		spectatorClose := state.SpectatorDone.NewListener()
		defer state.SpectatorDone.RemoveListener(spectatorClose)

		for {
			var buf []byte
			var err error

			select {
			case message := <-messages:
				buf = message.([]byte)

			case rawDoneMessage := <-spectatorClose:
				log.Println("[websocket] closing spectator socket")
				doneMessage, ok := rawDoneMessage.(*CloseMessage)

				var code int
				var text string
				if !ok {
					// Nil Message
					code = websocket.CloseNormalClosure
					text = ""
				} else {
					code = doneMessage.Code
					text = doneMessage.Text
				}

				// Teardown
				err = conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(code, text),
					time.Now().Add(time.Second),
				)
				if err != nil {
					log.Printf("[websocket] error while tearing down spectator: %v\n", err)
				}

				state.Destroy(session)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.Println("[websocket] failed to write", err)
				return
			}
		}
	}()

	defer state.SpectatorDone.Broadcast(nil)

	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("[game] spectator disconnected")
		state.SpectatorDone.Broadcast(&CloseMessage{
			Text: text,
			Code: code,
		})

		// Losing the error is ok. The connection continues to get torn down as
		// it normally does except any error from the close doesn't get passed
		// to the next call to Read() or NextReader(). A close error is still
		// sent in its place.

		// TODO: Change this. It relies on the semantics of the library
		// remaining the same.
		return nil
	})

	// Start is the only message spectators can send.
	err := ReadMessage(conn, dosProto.MessageType_START, nil)
	if err != nil {
		return
	}

	go func() {
		// Even though spectators aren't allowed to send anything, keep an
		// active read so the close handler works.
		err = Read(conn, nil)

		log.Println("[websocket] read data from spectator", err)

		// Client sent us some garbage. Time to take out the trash.
		state.SpectatorDone.Broadcast(&CloseMessage{
			Text: "",
			Code: websocket.CloseUnsupportedData,
		})

		return
	}()

	log.Println("[game] starting")

	state.IsStarted = true
	state.Started.Broadcast(nil)

	for {
		player := game.NextPlayer()
		if player == nil {
			// Game is done or there aren't enough players.
			return
		}

		log.Printf("[game] (%s) turn\n", player.Name)
		state.Turn.Broadcast(player.Name)

		select {
		// TODO: On Player Leave, If there aren't enough players, shutdown the
		// game.
		case <-player.TurnDone:
			log.Printf("[game] (%s) turn done\n", player.Name)
			if len(player.Cards.List) == 0 {
				log.Printf("[game] (%s) done with game\n", player.Name)

				state.PlayerDone.Broadcast(&CloseMessage{
					Name: player.Name,
					Code: websocket.CloseNormalClosure,
					Text: "won!",
				})
			}
		}
	}
}

func handlePlayer(conn *websocket.Conn) {
	log.Println("[websocket] new player joined")

	var state *GameState
	var attempts int
	for attempts = 0; state == nil && attempts < MaxAttempts; attempts++ {
		sessionMessage := &dosProto.SessionMessage{}
		err := ReadMessage(conn, dosProto.MessageType_SESSION, sessionMessage)
		if err != nil {
			return
		}

		var ok bool
		state, ok = GameStore[sessionMessage.Session]

		var message proto.Message
		var typ dosProto.MessageType
		if !ok {
			message = &dosProto.ErrorMessage{
				Reason: "Invalid Game PIN. That game doesn't exist.",
			}
			typ = dosProto.MessageType_ERROR
		} else {
			message = nil
			typ = dosProto.MessageType_SUCCESS
		}

		err = WriteMessage(conn, typ, message)
		if err != nil {
			return
		}
	}

	if attempts == MaxAttempts {
		// Ratelimited
		// We haven't setup any writing goroutines. We don't have anything to
		// teardown besides the connection so writing directly is ok here.
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "slow down"),
			time.Now().Add(time.Second),
		)
	}

	game := state.Game

	// Wait for player to be ready
	var player *dos.Player

	for attempts = 0; player == nil && attempts < MaxAttempts; attempts++ {
		ready := dosProto.ReadyMessage{}
		err := ReadMessage(conn, dosProto.MessageType_READY, &ready)
		if err != nil {
			return
		}

		var message proto.Message
		var typ dosProto.MessageType

		// TODO: possibly block empty player names ""
		player, err = game.NewPlayer(ready.Name)
		if err != nil {
			log.Printf("[game] (%s) failed to join: %v\n", ready.Name, err)

			typ = dosProto.MessageType_ERROR
			message = &dosProto.ErrorMessage{Reason: err.Error()}
		} else {
			log.Printf("[game] (%s) joined\n", ready.Name)

			typ = dosProto.MessageType_SUCCESS
			message = nil
		}

		err = WriteMessage(conn, typ, message)
		if err != nil {
			return
		}
	}

	if attempts == MaxAttempts {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "slow down"),
			time.Now().Add(time.Second),
		)
	}

	isStarted := false
	isMyTurn := false

	hasDrawn := false
	hasPlayed := false

	// This goroutine controls write access to the player's socket.
	go func() {
		// Send player list
		playersMessage := dosProto.PlayersMessage{Additions: game.GetPlayerList()}
		err := WriteMessage(conn, dosProto.MessageType_PLAYERS, &playersMessage)
		if err != nil {
			return
		}

		handAdditions := player.Cards.Additions.NewListener()
		go player.Cards.Additions.StartBroadcasting()
		defer player.Cards.Additions.Destroy()

		handDeletions := player.Cards.Deletions.NewListener()
		go player.Cards.Deletions.StartBroadcasting()
		defer player.Cards.Deletions.Destroy()

		messages := state.CommonMessages.NewListener()
		defer state.CommonMessages.RemoveListener(messages)

		start := state.Started.NewListener()
		defer state.Started.RemoveListener(start)

		turnChanged := state.Turn.NewListener()
		defer state.Turn.RemoveListener(turnChanged)

		playerDone := state.PlayerDone.NewListener()
		defer state.PlayerDone.RemoveListener(playerDone)

		spectatorDone := state.SpectatorDone.NewListener()
		defer state.SpectatorDone.RemoveListener(spectatorDone)

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

			case rawDoneMessage := <-playerDone:
				doneMessage := rawDoneMessage.(*CloseMessage)
				if doneMessage.Name != player.Name {
					// This message is for someone else.
					continue
				}

				// Teardown
				err = conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(doneMessage.Code, doneMessage.Text),
					time.Now().Add(time.Second),
				)

				if err != nil {
					log.Println("[websocket] error while tearing down player socket:", err)
				}

				game.RemovePlayer(player)

				// Notify Spectator If Turn Active, Otherwise Ignored
				player.TurnDone <- struct{}{}

				return

			case <-spectatorDone:
				// Teardown Connection. Everything Else Will Be GC'd
				err = conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
					time.Now().Add(time.Second),
				)

				if err != nil {
					log.Println("[websocket] error occured while spectator closing player connection:", err)
				}

				return
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

	// Handle leaving
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("[game] (%s) is leaving\n", player.Name)
		state.PlayerDone.Broadcast(&CloseMessage{
			Name: player.Name,
			Code: code,
			Text: text,
		})

		return nil
	})

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

		default:
			state.PlayerDone.Broadcast(&CloseMessage{
				Text: "you're using that wrong",
				Code: websocket.CloseUnsupportedData,
			})
			return
		}
	}
}
