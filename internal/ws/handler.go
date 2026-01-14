package ws

import (
	"go-push-service/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		logger.Warn("Connection attempt without user_id rejected")
		http.Error(w, "Unauthorized: user_id required", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err, "WebSocket upgrade failed")
		return
	}

	client := &Client{
		Hub:    hub,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Rdb:    rdb,
	}

	client.Hub.Register <- client

	logger.Infof("user_id", userID, "Client connected")

	go client.WritePump()
	go client.ReadPump()
	go client.RedisPump()
}
