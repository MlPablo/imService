package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	"imService/broker"
	"imService/storage"
)

// server is struct for server
type server struct {
	router *gin.Engine
	store  storage.Storage
	broker broker.Broker
}

// NewServer creates server, setting: routes storage, producer and starts the producer and return the server
func NewServer() *server {
	if err := gotenv.Load(".env"); err != nil {
		panic(err)
	}
	s := &server{
		router: gin.Default(),
		store:  storage.NewStorage(),
		broker: broker.NewQue(),
	}

	s.router.LoadHTMLGlob("templates/*.html")
	s.setRoutes()
	go s.broker.Consume(s.store)

	return s
}

// setRoutes set routes with handlers
func (s *server) setRoutes() {
	s.router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})
	s.router.POST("/upload", SaveFile(s.broker, s.store))
	s.router.GET("/download/:id", GetFiles(s.store))
}

// Start starts server on port 8080
func Start() {
	server := NewServer()
	log.Fatal(server.router.Run(":8080"))
}
