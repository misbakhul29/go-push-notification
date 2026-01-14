package main

import (
	"go-push-service/internal/config"
	"go-push-service/internal/mq"
	"go-push-service/internal/store"
	"go-push-service/internal/ws"
	"go-push-service/pkg/logger"
	"net/http"
)

func main() {
	logger.Setup()

	cfg := config.LoadConfig()
	logger.Infof("config", cfg, "Configuration loaded")

	rdb := store.NewRedisClient(cfg.RedisAddr, cfg.RedisPass)
	logger.Info("Redis connection initialized")

	hub := ws.NewHub()
	go hub.Run()
	logger.Info("WebSocket Hub running")

	mq.StartConsumer(cfg, rdb)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, rdb, w, r)
	})

	logger.Info("Server started on port " + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Fatal(err, "Server crashed")
	}
}
