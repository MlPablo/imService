package server

import (
	"github.com/gin-gonic/gin"
	"imService/rabbit"
	"imService/storage"
	"net/http"
)

// Just using static Queue name and path to connect to RabbitMQ (commented path is using for docker-compose)
const (
	que = "ImageQue"
	//path = "amqp://guest:guest@rabbitmq/"
	path = "amqp://guest:guest@localhost:5672/"
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
	s := &server{
		router:   gin.Default(),
		store:    storage.NewStorage(),
		queue:    rabbit.NewProducer(que, path),
		consumer: rabbit.NewConsumer(que, path),
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
	server.router.Run(":8080")
}
