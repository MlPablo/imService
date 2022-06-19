package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
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

// CreateFileInDownloads Create a new file in which image will be copied (dir for file is ...\Downloads\)
func CreateFileInDownloads(id, quality, ext string) (*os.File, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	fmt.Println(dirname)
	file, err := os.Create(fmt.Sprintf("%s\\Downloads\\id_%s_%s.%s", dirname, id, quality, ext))
	if err != nil {
		return nil, err
	}
	return file, nil
}
