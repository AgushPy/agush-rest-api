package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"curso-rest.com/go/rest/database"
	"curso-rest.com/go/rest/repository"
	"curso-rest.com/go/rest/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	//Puerto donde se va a ejecutar
	Port string
	//Clave secret para poder generar tokens
	JWTSecret string
	//Conexión a la base de datos
	DatabaseURL string
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
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}
	if config.DatabaseURL == "" {
		return nil, errors.New("database is required")
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
	//Como estoy implementando la interface de Server lo puedo enviar al broker como un server
	binder(b, b.router)
	// handler := cors.Default().Handler(b.router)
	handler := cors.AllowAll().Handler(b.router)
	repo, err := database.NewMysqlRepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting Server on port", b.Config().Port)

	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServer: ", err)
	} else {
		log.Fatal("server stopped")
	}

}
