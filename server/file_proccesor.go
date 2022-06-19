package server

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"strings"
)

// FileProcessor accepts a file that was uploaded, check it for right extension and remakes it into bytes and returns
func FileProcessor(header *multipart.FileHeader) ([]byte, error) {
	ext := strings.Split(strings.ToLower(header.Filename), ".")[1]
	if ext != "png" && ext != "jpg" {
		return nil, errors.New("should be png or jpeg format")
	}

	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
