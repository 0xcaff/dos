package main

import (
	"fmt"
	"github.com/caffinatedmonkey/dos/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var wg = sync.WaitGroup{}
var upgrader = websocket.Upgrader{}
var connections = map[string]*websocket.Conn{}
var spectatorConnections = []*websocket.Conn{}

var mux = http.NewServeMux()
var s = &http.Server{
	Addr:    ":8080",
	Handler: mux,
}
var game = dos.NewGame()

// TODO: Put these in a different file
type startGameRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type turnResponse struct {
	Type     string `json:"type"`
	dos.Card `json:"card"`
}

func main() {
	mux.HandleFunc("/ws/play", handlePlay)
	mux.HandleFunc("/ws/spectate", handleSpectate)
	mux.Handle("/", http.FileServer(http.Dir("assets")))

	wg.Add(1)
	go func(s *http.Server) {
		fmt.Println("Starting Server")
		log.Fatal(s.ListenAndServe())
		wg.Done()
	}(s)

	wg.Wait()
}

func handlePlay(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		// TODO: Handle Error
		log.Fatal(err)
		panic(err)
	}

	sgr := &startGameRequest{}
	err = conn.ReadJSON(sgr)
	if err != nil {
		// TODO: Handle Error
		log.Fatal(err)
		panic(err)
	}

	if sgr.Type == "start" {
		// Add Players to the Game
		player := game.NewPlayer(sgr.Name)
		fmt.Printf("Player (%s) Joined the Game!\n", player.Name)
		connections[player.Name] = conn

		for _, specConn := range spectatorConnections {
			player.Hand.OnAddition(func(i interface{}) {
				specConn.WriteJSON(map[string]interface{}{
					"type":  "addition",
					"to":    "hand",
					"name":  player.Name,
					"added": i,
				})
			})
		}
	}
}

func handleSpectate(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		// TODO: Handle Error
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Spectator Joined")
	conn.WriteJSON(map[string]interface{}{
		"type":     "initialInformation",
		"deck":     game.Deck.Cards.List,
		"players":  game.Players,
		"lastCard": game.LastPlayedCard,
	})

	game.Deck.Cards.OnDeletion(func(i interface{}) {
		conn.WriteJSON(map[string]interface{}{
			"type":    "deletion",
			"from":    "deck",
			"removed": i,
		})
	})

	game.Players.OnAddition(func(i interface{}) {
		conn.WriteJSON(map[string]interface{}{
			"type":  "addition",
			"to":    "player",
			"added": i,
		})
	})

	game.Players.OnDeletion(func(i interface{}) {
		conn.WriteJSON(map[string]interface{}{
			"type":    "deletion",
			"from":    "player",
			"removed": i,
		})
	})
	spectatorConnections = append(spectatorConnections, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		if messageType == websocket.TextMessage && string(p) == "start" {
			fmt.Println("Starting Game")

			// Send player information to other players
			for _, p := range game.Players.List {
				player := p.(dos.Player)
				conn := connections[player.Name]
				err := conn.WriteJSON(map[string]interface{}{
					"type": "start",
					// TODO: Player Names only
					"players": game.Players,
					"hand":    player.Hand.List,
				})
				if err != nil {
					log.Fatal(err)
				}
			}

			// Game Loop
			for {
				player := game.NextPlayer()
				playerConn := connections[player.Name]

				for _, p := range game.Players.List {
					player := p.(dos.Player)
					conn := connections[player.Name]
					err := conn.WriteJSON(map[string]interface{}{
						"type": "updateCard",
						"card": game.LastPlayedCard,
					})
					if err != nil {
						log.Fatal(err)
					}
				}

				fmt.Printf("Player (%s) Turn\n", player.Name)
				err = playerConn.WriteJSON(map[string]interface{}{
					"type": "turn",
				})
				if err != nil {
					log.Fatal(err)
				}

				for {
					tr := turnResponse{}
					err = playerConn.ReadJSON(&tr)
					if err != nil {
						log.Fatal(err)
					}

					// TODO: fix ability to draw multiple cards per turn
					if tr.Type == "drawCard" {
						newCard := game.Deck.Cards.Pop()
						player.Hand.Push(newCard)
						fmt.Printf("Player (%s), Drew Card, %v\n", player.Name, newCard)
					} else if tr.Type == "playCard" {
						game.PlayCard(tr.Card)
						fmt.Printf("Player (%s), Played Card %v\n", player.Name, tr.Card)
						break
					}

					// TODO: End Game Thingy
					// TODO: End of DrawDeck Reshuffler Thing
				}
			}
		}
	}
}
