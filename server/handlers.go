package server

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"imService/rabbit"
	"imService/storage"
)

// Set image size limit to upload
const maxBytes = 1024 * 1024

// SaveFile get an uploaded image, checks and processes it. If everything ok, adds file to queue by rabbit.RMQProducer
func SaveFile(q rabbit.RMQProducer, db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
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

		filename := fmt.Sprintf("id_%s_%s.%s", id, quality, strings.Split(imag.Ext, "/")[1])
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
		c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))
		c.Writer.Header().Set("Content-Length", c.Request.Header.Get("Content-Length"))

		switch imag.Ext {
		case "image/jpeg":
			err = jpeg.Encode(c.Writer, imag.Image, nil)
		case "image/png":
			err = png.Encode(c.Writer, imag.Image)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Successfully downloaded"})
	}
}
