package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func (app *Application) StartServer(port int) {
	// HTTP endpoint to fetch UUIDs
	http.HandleFunc("/g", UUIDHandler(app))

	// WebSocket endpoint to receive UUIDs
	http.HandleFunc("/ws/g", WebSocketHandler(app))

	// Start HTTP server
	fmt.Println("serverd")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("HTTP server error: %v\n", err)
	}
}

func UUIDHandler(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nStr := r.URL.Query().Get("n")
		n, err := strconv.Atoi(nStr)
		if err != nil || n <= 0 {
			n = 1
		}
		if n > 1000 {
			n = 1000
		}

		uuids := app.GetUUIDs(n)

		jsonResponse, err := json.Marshal(uuids)
		if err != nil {
			http.Error(w, "Error generating UUIDs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

type wsMessage struct {
	N int `json:"n"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}

		go handleMessages(conn, app)
	}
}

func handleMessages(conn *websocket.Conn, app *Application) {
	defer conn.Close()

	for {
		var msg wsMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			return
		}

		n := msg.N
		if n <= 0 {
			n = 1
		}
		if n > 1000 {
			n = 1000
		}

		uuids := app.GetUUIDs(n)

		if err := conn.WriteJSON(uuids); err != nil {
			return
		}
	}
}
