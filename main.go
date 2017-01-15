//go:generate protoc --go_out=. proto/card.proto proto/handshake.proto proto/ready.proto proto/changed.proto proto/players.proto proto/turn.proto proto/envelope.proto proto/play.proto proto/session.proto
package main

import (
	"flag"
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
	log "github.com/sirupsen/logrus"
)

// Maximum failed attempts to join a session or select a name before the
// connection is terminated.
const MaxAttempts = 10

var GameStore = make(map[int32]*GameState)
var StoreMutex sync.RWMutex

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

	log.Infof("initializing server at %s", *listen)

	err := s.ListenAndServe()
	log.Println(err)
}

func handleSocket(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Warning("connection initialization failed", err)
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

	default:
		// Close the connection. This is an invalid type.
		log.Warning("invalid connection type")

		err = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "invalid client type"),
			time.Now().Add(time.Second),
		)
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
			log.Error("protobuf encoding error", err)
		} else {
			state.CommonMessages.Broadcast(bytes)
		}
	}
}

// Removes the state from the global store.
// TODO: All broadcasters started in NewGame need to be destroyed.
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

func (closeMessage *CloseMessage) AsCloseError() *websocket.CloseError {
	return &websocket.CloseError{
		Code: closeMessage.Code,
		Text: closeMessage.Text,
	}
}

func handleSpectator(conn *websocket.Conn) {
	log.Debug("spectator connection")

	state, session := NewGame()
	logCtx := log.Fields{
		"session": session,
	}
	log.WithFields(logCtx).Info("new game")

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
				log.WithFields(logCtx).Info("closing spectator socket")
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
					log.WithFields(logCtx).Error("error while tearing down spectator socket: ", err)
				}

				conn.Close()
				state.Destroy(session)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.WithFields(logCtx).Error("failed to write to spectator socket", err)
				return
			}
		}
	}()

	defer state.SpectatorDone.Broadcast(nil)

	conn.SetCloseHandler(func(code int, text string) error {
		log.WithFields(logCtx).Info("spectator disconnecting")
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
		if err != nil {
			// A close message or protocol error. Either way the connection is
			// done.
			return
		}

		// The client read was successful. It should only be closed from this
		// point forward.
		state.SpectatorDone.Broadcast(&CloseMessage{
			Text: "",
			Code: websocket.CloseUnsupportedData,
		})

		return
	}()

	log.WithFields(logCtx).Info("starting")

	state.IsStarted = true
	state.Started.Broadcast(nil)

	for {
		player := game.NextPlayer()
		if player == nil {
			// Game is done or there aren't enough players.
			return
		}

		state.Turn.Broadcast(player.Name)

		// Wait until player is done with their turn.
		<-player.TurnDone
	}
}

func handlePlayer(conn *websocket.Conn) {
	log.Debug("player connection")
	logCtx := log.Fields{}

	var state *GameState

	var attempts int
	for attempts = 0; state == nil && attempts < MaxAttempts; attempts++ {
		sessionMessage := &dosProto.SessionMessage{}
		err := ReadMessage(conn, dosProto.MessageType_SESSION, sessionMessage)
		if err != nil {
			return
		}

		var ok bool
		StoreMutex.RLock()
		state, ok = GameStore[sessionMessage.Session]
		StoreMutex.RUnlock()

		var message proto.Message
		var typ dosProto.MessageType
		if !ok {
			message = &dosProto.ErrorMessage{Reason: dosProto.ErrorReason_INVALIDGAME}
			typ = dosProto.MessageType_ERROR
		} else if state.IsStarted {
			message = &dosProto.ErrorMessage{Reason: dosProto.ErrorReason_GAMESTARTED}
			typ = dosProto.MessageType_ERROR
		} else {
			message = nil
			typ = dosProto.MessageType_SUCCESS

			logCtx["session"] = sessionMessage.Session
		}

		err = WriteMessage(conn, typ, message)
		if err != nil {
			return
		}
	}

	if attempts == MaxAttempts {
		logCtx["remoteAddr"] = conn.RemoteAddr()
		log.WithFields(logCtx).Warning("rate limited while sending game PIN")

		// We haven't setup any writing goroutines. We don't have anything to
		// teardown besides the connection so writing directly is ok here.
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "slow down"),
			time.Now().Add(time.Second),
		)

		conn.Close()
		return
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

		logCtx["name"] = ready.Name

		var message proto.Message
		var typ dosProto.MessageType

		player, err = game.NewPlayer(ready.Name)
		if err != nil {
			log.WithFields(logCtx).Info("failed to join", err)

			// Only error we can get is invalid name.
			typ = dosProto.MessageType_ERROR
			message = &dosProto.ErrorMessage{Reason: dosProto.ErrorReason_INVALIDNAME}
		} else {
			log.WithFields(logCtx).Info("joined")

			typ = dosProto.MessageType_SUCCESS
			message = nil
		}

		err = WriteMessage(conn, typ, message)
		if err != nil {
			return
		}
	}

	if attempts == MaxAttempts {
		logCtx["remoteAddr"] = conn.RemoteAddr()
		log.WithFields(logCtx).Warning("rate limited while sending name")

		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "slow down"),
			time.Now().Add(time.Second),
		)

		conn.Close()
		return
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
				if isMyTurn {
					log.WithFields(logCtx).Info("turn")
					hasDrawn = false
					hasPlayed = false
				}
				continue

			case rawDoneMessage := <-playerDone:
				doneMessage := rawDoneMessage.(*CloseMessage)
				if doneMessage.Name != player.Name {
					// This message is for someone else.
					continue
				}
				log.WithFields(logCtx).Info("leaving ", doneMessage.AsCloseError())

				// Teardown
				err = conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(doneMessage.Code, doneMessage.Text),
					time.Now().Add(time.Second),
				)

				if err != nil {
					log.WithFields(logCtx).Error("error while closing player socket:", err)
				}

				// TODO: This might be closing in the wrong state. If the close
				// request is server initiated, we should wait until the client
				// responds back.
				conn.Close()
				game.RemovePlayer(player)

				// Notify Spectator
				player.TurnDone <- struct{}{}

				return

			case <-spectatorDone:
				// Teardown Connection. Everything Else Will Be Handled By
				// Defered
				err = conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
					time.Now().Add(time.Second),
				)

				if err != nil {
					log.WithFields(logCtx).Error("error occured while spectator closing player socket:", err)
				}
				conn.Close()

				return
			}

			if err != nil {
				log.WithFields(logCtx).Error("protobuf encoding error", err)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.WithFields(logCtx).Error("failed to write message", err)
				return
			}
		}
	}()

	// Handle leaving
	conn.SetCloseHandler(func(code int, text string) error {
		state.PlayerDone.Broadcast(&CloseMessage{
			Name: player.Name,
			Code: code,
			Text: text,
		})

		return nil
	})

	var closeMessage string
	var closeStatus int

	defer func() {
		// Since this is a nop and can be called multiple times, remove the
		// player. This needs to happen before TurnDone emits a message.
		game.RemovePlayer(player)

		if closeMessage != "" {
			// TODO: Race-y.
			// The sending of this closing info may not happen if the spectator
			// terminates the connection before the player does.
			state.PlayerDone.Broadcast(&CloseMessage{
				Name: player.Name,
				Code: closeStatus,
				Text: closeMessage,
			})
		}

		// Continues to the next player if it is our turn otherwise does
		// nothing.
		player.TurnDone <- struct{}{}
	}()

	for {
		envelope := dosProto.Envelope{}
		err := Read(conn, &envelope)
		if err != nil {
			// This error is because the connection closed.
			log.WithFields(logCtx).Debug("failed to read message: ", err)
			return
		}

		if !isMyTurn {
			log.WithFields(logCtx).Warning("tried to play during other person's turn")
			closeMessage = "you can't do that now"
			closeStatus = websocket.ClosePolicyViolation
			return
		}

		switch envelope.Type {
		case dosProto.MessageType_DRAW:
			if hasDrawn || hasPlayed {
				log.WithFields(logCtx).Warning("trying to draw after drawing or playing")
				closeMessage = "you can't do that now"
				closeStatus = websocket.ClosePolicyViolation
				return
			}

			log.WithFields(logCtx).Info("drawing card")

			game.DrawCards(&player.Cards, 1)
			hasDrawn = true

		case dosProto.MessageType_PLAY:
			if hasPlayed {
				log.WithFields(logCtx).Warning("trying to play multiple times")
				closeMessage = "you can't do that now"
				closeStatus = websocket.ClosePolicyViolation
				return
			}

			playMessage := dosProto.PlayMessage{}
			err := proto.Unmarshal(envelope.Contents, &playMessage)
			if err != nil {
				log.WithFields(logCtx).Warning("decoding error", err)
				closeMessage = "couldn't parse message"
				closeStatus = websocket.ClosePolicyViolation
				return
			}

			err = game.PlayCard(player, playMessage.Id, playMessage.Color)
			if err != nil {
				log.WithFields(logCtx).Warning("failed to play card", err, playMessage)
				closeMessage = "you can't play that card"
				closeStatus = websocket.ClosePolicyViolation
				return
			}

			log.WithFields(logCtx).Info("played card")
			hasPlayed = true

		case dosProto.MessageType_DONE:
			if !hasDrawn && !hasPlayed {
				log.WithFields(logCtx).Warning("trying to skip turn")
				closeMessage = "you can't do that now"
				closeStatus = websocket.ClosePolicyViolation
				return
			}

			if len(player.Cards.List) == 0 {
				log.WithFields(logCtx).Info("won")
				closeMessage = "won!"
				closeStatus = websocket.CloseNormalClosure
				return
			}
			player.TurnDone <- struct{}{}

		default:
			log.WithFields(logCtx).Warning("sent invalid message type")
			closeMessage = "you're using that wrong"
			closeStatus = websocket.ClosePolicyViolation
			return
		}
	}
}
