package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	dosProto "github.com/0xcaff/dos/proto"
	"github.com/0xcaff/dos/utils"
	"github.com/golang/protobuf/proto"
)

func TestWastedBroadcasters(t *testing.T) {
	// Setup Server
	server := httptest.NewServer(http.HandlerFunc(handleSocket))
	defer server.Close()

	// Rewrite URL
	url := strings.Replace(server.URL, "http", "ws", 1)

	// Connect To Server
	spectatorConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Send Spectator Handshake
	handshake := dosProto.HandshakeMessage{Type: dosProto.ClientType_SPECTATOR}
	err = Write(spectatorConn, &handshake)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Get Game ID
	sessionMessage := &dosProto.SessionMessage{}
	err = ReadMessage(spectatorConn, dosProto.MessageType_SESSION, sessionMessage)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	go IgnoreIncoming(spectatorConn)

	// Make Sure Game Exists
	gameState, ok := GameStore[sessionMessage.Session]
	if !ok {
		t.Log("The game we got back doesn't exist")
		t.Fail()
	}

	// Make Sure That's The Only Game
	if len(GameStore) > 1 {
		t.Log("More than one game in the game store.")
		t.Fail()
	}

	commonMessages := gameState.CommonMessages.NewListener()
	player1Conn, err := AddPlayer(url, "Player 1", sessionMessage.Session)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	go IgnoreIncoming(player1Conn)
	<-commonMessages

	player2Conn, err := AddPlayer(url, "Player 2", sessionMessage.Session)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	go IgnoreIncoming(player2Conn)
	<-commonMessages

	// Start Game
	err = WriteMessage(spectatorConn, dosProto.MessageType_START, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	<-commonMessages
	gameState.CommonMessages.RemoveListener(commonMessages)

	// Close Connection
	err = spectatorConn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseGoingAway, ""),
		time.Now().Add(time.Second),
	)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// In this case all connection closes come from the client.
	spectatorConn.SetCloseHandler(func(code int, text string) error {
		// We got a message giving the go ahead to close the connection.

		err = spectatorConn.Close()
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		return nil
	})

	// TODO: Wait until the end of the shutdown process. This waits until near
	// the end but other goroutines could still be shutting down.
	for {
		StoreMutex.RLock()
		_, ok := GameStore[sessionMessage.Session]
		StoreMutex.RUnlock()
		if !ok {
			break
		}
	}

	// TODO: With even this much time, the shutdown never completes.
	time.Sleep(time.Second * 5)

	// Count Wasted Broadcasters
	toCheckIfClosed := []utils.Broadcaster{
		gameState.CommonMessages,
		gameState.Started,
		gameState.Turn,
		gameState.PlayerDone,
		gameState.SpectatorDone,

		gameState.Game.Deck.Additions,
		gameState.Game.Deck.Deletions,

		gameState.Game.Discard.Additions,
		gameState.Game.Discard.Deletions,
	}

	// Check Player Card Collections Too. The players aren't removed from the
	// game since the spectator initiated the shutdown.
	players := gameState.Game.GetPlayers()
	for _, player := range players {
		toCheckIfClosed = append(toCheckIfClosed, player.Cards.Additions, player.Cards.Deletions)
	}

	for i, broadcaster := range toCheckIfClosed {
		if !broadcaster.IsClosed {
			t.Log("Wasted Broadcaster", i)
			t.Log("With Listeners", broadcaster.CountListeners())
			t.Fail()
		}
	}
}

func AddPlayer(url, name string, gameId int32) (*websocket.Conn, error) {
	// Connect To Server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	// Send Player Handshake
	handshake := dosProto.HandshakeMessage{Type: dosProto.ClientType_PLAYER}
	err = Write(conn, &handshake)
	if err != nil {
		return nil, err
	}

	// Join Game
	sessionMessage := dosProto.SessionMessage{Session: gameId}
	err = WriteMessage(conn, dosProto.MessageType_SESSION, &sessionMessage)
	if err != nil {
		return nil, err
	}

	// Send Name
	readyMessage := dosProto.ReadyMessage{Name: name}
	err = WriteMessage(conn, dosProto.MessageType_READY, &readyMessage)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func IgnoreIncoming(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func Write(conn *websocket.Conn, message proto.Message) error {
	buf, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Test Rate Limiting and Connection Rejection
