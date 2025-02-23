package util

import (
	"fmt"
	"github.com/google/uuid"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func DecodeImage(imageBytes []byte) (filepath string, err error) {
	id := uuid.New().String()
	name := fmt.Sprintf("%s.png", id)
	file, err := os.Create("./files/" + name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(imageBytes)
	if err != nil {
		return "", err
	}
	filepath = "./files/" + name
	err = nil
	return
}
