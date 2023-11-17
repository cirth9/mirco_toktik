package main

import (
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
)

func saveFile(header *multipart.FileHeader, destination string) error {
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer func(file multipart.File) {
		err1 := file.Close()
		if err1 != nil {
			zap.S().Error(err1)
		}
	}(file)
	newFile, err := os.Create(destination)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	written, err := io.Copy(newFile, file)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	zap.S().Info("written ", written, " to ", destination)
	return nil
}

func main() {
	_, err := os.Create("../static/123.mp4")
	if err != nil {
		zap.S().Error(err)
	}
}
