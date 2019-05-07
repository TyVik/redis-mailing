package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

func receive(ctx context.Context, ps *redis.PubSub, consumer chan<- *redis.Message) {
	ch := ps.Channel()
	defer ps.Close()

	for {
		select {
		case msg := <-ch:
			consumer <- msg
		case <-ctx.Done():
			return
		}
	}
}

func consume(ctx context.Context, conn *websocket.Conn, ch <-chan *redis.Message, log zerolog.Logger) {
	for {
		select {
		case msg := <-ch:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.String())); err != nil {
				log.Error().Err(err).Msg("write msg")
			}
		case <-ctx.Done():
			return
		}
	}
}

func runSever(cancel context.CancelFunc, addr string, handler http.Handler, log zerolog.Logger) {
	defer cancel()

	if err := http.ListenAndServe(addr, handler); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("failed main server")
	}
}

func main() {
	cfg := &config{}
	envconfig.MustProcess("", cfg)

	log := getLogger(cfg.LogLevel)

	bus := make(chan *redis.Message)
	defer close(bus)

	redisClinet := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClinet.Close()

	if err := redisClinet.Ping().Err(); err != nil {
		log.Fatal().Err(err).Msg("connect to redis")
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		receive(ctx, redisClinet.Subscribe(cfg.PubSubChannels...), bus)
	}()

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("upgrade conn")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		consume(ctx, conn, bus, log)
	})

	runSever(cancel, ":"+cfg.ServicePort, accessHandler(log, router), log)

	wg.Wait()
}
