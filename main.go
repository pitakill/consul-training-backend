package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const VERSION = "3.0.0"
const key = "key"

func logic(client *redis.Client) (int, error) {
	ctx := context.Background()
	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		const value = 1
		return value, client.Set(ctx, key, value, 0).Err()
	} else if err != nil {
		return 0, err
	}

	current, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	actual := current + 1

	return actual, client.Set(ctx, key, actual, 0).Err()
}

func healthcheck(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	log.Println(r)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(struct {
		Success string `json:"success"`
	}{"ok"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := w.Write(j); err != nil {
		log.Fatal(err)
	}
}

func root(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	log.Println(r)
	// App logic (very basic)
	count, err := logic(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle response
	j, err := json.Marshal(struct {
		Count   int    `json:"count"`
		Version string `json:"version"`
	}{count, VERSION})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := w.Write(j); err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	redis_host := os.Getenv("REDIS_HOST")
	if redis_host == "" {
		redis_host = "localhost"
	}
	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		redis_port = "6379"
	}
	redis_password := os.Getenv("REDIS_PASSWORD")
	if redis_password == "" {
		redis_password = "holamundo"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
		Password: redis_password,
		DB:       0,
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) { healthcheck(w, r, rdb) })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { root(w, r, rdb) })

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
