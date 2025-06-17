package logstream

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type    string      `json:"type"`    // "log" or "job"
	Content string      `json:"content"` // for logs
	Job     interface{} `json:"job"`     // for job updates
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	// Send initial jobs if needed
	// You might want to add this after implementing store.GetAllJobs()
	// initialJobs := store.GetAllJobs()
	// if err := conn.WriteJSON(Message{Type: "init", Job: initialJobs}); err != nil {
	//     return
	// }

	for {
		// Just keep the connection alive
		if _, _, err := conn.NextReader(); err != nil {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			break
		}
	}
}

func BroadcastLog(msg string) {
	message := Message{
		Type:    "log",
		Content: msg,
	}
	broadcastMessage(message)
}

func BroadcastJob(job interface{}) {
	message := Message{
		Type: "job",
		Job:  job,
	}
	broadcastMessage(message)
}

func broadcastMessage(msg Message) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}
