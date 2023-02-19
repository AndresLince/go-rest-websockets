package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/AndresLince/go-rest-websockets/database"
	"github.com/AndresLince/go-rest-websockets/repository"
	"github.com/AndresLince/go-rest-websockets/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Addr        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Addr == "" {
		return nil, errors.New("addr is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("database Url is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	handler := cors.Default().Handler(b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on addr", b.config.Addr)
	err = http.ListenAndServe(b.config.Addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
