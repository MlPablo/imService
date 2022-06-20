package server

import (
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"imService/rabbit"
	"imService/storage"
	"log"
	"net/http"
	"os"
)

// server is struct for server
type server struct {
	router   *gin.Engine
	store    storage.Storage
	queue    rabbit.RMQProducer
	consumer rabbit.RMQConsumer
}

// NewServer creates server, setting: routes storage, producer and starts the producer and return the server
func NewServer() *server {
	gotenv.Load(".env")
	s := &server{
		router:   gin.Default(),
		store:    storage.NewStorage(),
		queue:    rabbit.NewProducer(os.Getenv("RABBIT_QUE"), os.Getenv("RABBIT_PATH")),
		consumer: rabbit.NewConsumer(os.Getenv("RABBIT_QUE"), os.Getenv("RABBIT_PATH")),
	}

	s.router.LoadHTMLGlob("templates/*.html")
	s.setRoutes()
	go s.consumer.Consume(s.store)

	return s
}

// setRoutes set routes with handlers
func (s *server) setRoutes() {
	s.router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})
	s.router.POST("/upload", SaveFile(s.queue, s.store))
	s.router.GET("/download/:id", GetFiles(s.store))
}

// Start starts server on port 8080
func Start() {
	server := NewServer()
	log.Fatal(server.router.Run(":8080"))
}
