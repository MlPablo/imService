package server

import (
	"github.com/gin-gonic/gin"
	"imService/rabbit"
	"imService/storage"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

// SaveFile get an uploaded image, checks and processes it. If everything ok, adds file to queue by rabbit.RMQProducer
func SaveFile(q rabbit.RMQProducer, db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var maxBytes int64 = 1024 * 1024 // 1MB
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, maxBytes)

		handler, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		bytesImage, err := FileProcessor(handler)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		q.PublishMessage(http.DetectContentType(bytesImage), bytesImage)

		c.JSON(200, gin.H{"message": "Saved! Your id is " + db.GetCurrentId()})
	}
}

// GetFiles creates a file at ...\Downloads\ and copying into it a requested image if it exists in DataBase
func GetFiles(db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		quality := c.Query("quality")

		imag, err := db.Get(id, quality)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		file, err := CreateFileInDownloads(id, quality, strings.Split(imag.Ext, "/")[1])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer file.Close()

		switch imag.Ext {
		case "image/jpeg":
			err = jpeg.Encode(file, imag.Image, nil)
		case "image/png":
			err = png.Encode(file, imag.Image)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Successfully downloaded"})
	}
}
