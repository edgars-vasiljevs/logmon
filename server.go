package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

const HOST string = "127.0.0.1"
const PORT string = "8080"

var clients []*websocket.Conn

var publicFiles = map[string]string{
	"assets/index.html": "/",
}

type LogPreload struct {
	Files [][]string `json:"files"`
}

// Listen for websocket connections and broadcast logs
func NewWebSocketServer(config Config, logs <-chan LogMessage) {

	endChan := make(chan *websocket.Conn)

	items := [][]string{}
	for _, item := range config.content {
		items = append(items, item)
	}

	// Listen to incoming connections
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {

		// Send list of files first
		websocket.JSON.Send(ws, LogPreload{items})

		// Add client to client list
		clients = append(clients, ws)

		// Wait for lost connection
		for deadWs := range endChan {
			if deadWs == ws { break }
		}
	}))

	// Broadcast logs to all clients
	for log := range logs {
		for i, client := range clients {
			err := websocket.JSON.Send(client, log)

			// Delete dead connection
			if err != nil {
				clients = append(clients[:i], clients[i+1:]...)
				endChan <- client
			}
		}
	}
}

// Serve public files over HTTP for front-end UI
func NewHTTPServer() {

	for file, path := range publicFiles {
		http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			assetContent, _ := Asset(file)
			fmt.Fprint(w, strings.Replace(string(assetContent), "{port}", PORT, -1))
		})
	}

	Print("Listening on http://" + HOST + ":" + PORT + "/")
	err := http.ListenAndServe(HOST + ":" + PORT, nil)

	if err != nil {
		Print(err)
	}
}
