package server_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"imService/rabbit"
	"imService/server"
	"imService/storage"
	"imService/test_image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func EmptyTestQueue(db storage.Storage) {
	consumer := rabbit.NewConsumer("TestImageQue", "amqp://guest:guest@localhost:5672")
	go consumer.Consume(db)
}

func TestServer_UploadImage(t *testing.T) {
	store := storage.NewStorage()
	EmptyTestQueue(store)
	queue := rabbit.NewProducer("TestImageQue", "amqp://guest:guest@localhost:5672")
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	r := SetUpRouter()
	r.POST("/upload", server.SaveFile(queue, store))

	testcases := []struct {
		name         string
		expectedCode int
		file         string
	}{
		{
			name:         "valid file",
			expectedCode: 200,
			file:         "ILTQq.png",
		},
		{
			name:         "invalid file type",
			expectedCode: 400,
			file:         "main.go",
		},
		{
			name:         "invalid file type",
			expectedCode: 400,
			file:         "abk.dd",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			go func() {
				defer writer.Close()
				_, err := writer.CreateFormFile("file", tc.file)
				if err != nil {
					t.Error(err)
				}
			}()
			req, _ := http.NewRequest("POST", "/upload", pr)
			req.Header.Add("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tc.expectedCode)
		})
	}
}

func TestServer_GetFile(t *testing.T) {
	store := storage.NewStorage()
	file, _ := os.Open(test_image.PathToTestImage)
	image, _ := png.Decode(file)
	err := store.Add(image, "image/png")
	t.Error(err)
	r := SetUpRouter()
	r.GET("/download/:id", server.GetFiles(store))

	testcases := []struct {
		name         string
		id           string
		quality      string
		downloaded   bool
		expectedCode int
	}{
		{
			name:         "valid request",
			id:           "1",
			quality:      "50",
			downloaded:   true,
			expectedCode: 200,
		},
		{
			name:         "invalid request id",
			id:           "10",
			quality:      "25",
			downloaded:   false,
			expectedCode: 400,
		},
		{
			name:         "invalid request quality",
			id:           "1",
			quality:      "28",
			downloaded:   false,
			expectedCode: 400,
		},
		{
			name:         "invalid request id and quality",
			id:           "e5",
			quality:      "fvt",
			downloaded:   false,
			expectedCode: 400,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/download/%s", tc.id), nil)
			q := req.URL.Query()
			q.Add("quality", tc.quality)
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedCode)
			dir, _ := os.UserHomeDir()
			_, err := os.Open(fmt.Sprintf("%s\\Downloads\\id_%s_%s.%s", dir, tc.id, tc.quality, "png"))
			if tc.downloaded {
				assert.Equal(t, err, nil)
			} else {
				assert.NotEqual(t, err, nil)
			}

		})
	}

}
