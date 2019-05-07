package main

import (
	"bytes"
	"log"
	"net/url"
	"sync"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		redisClinet := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		defer redisClinet.Close()

		if err := redisClinet.Ping().Err(); err != nil {
			log.Fatalf("ping redis: %v", err)
		}

		err := redisClinet.Publish("testChannel", "hello").Err()
		if err != nil {
			log.Fatalf("publish msg: %v", err)
		}
	}()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error dial:", err)
	}
	defer c.Close()

	wg.Wait()

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatalf("error read msg: %v", err)
	}

	if !bytes.Equal(message, []byte("Message<testChannel: hello>")) {
		log.Fatalf("something wrong: %s", string(message))
	} else {
		log.Println("OK!")
	}
}
