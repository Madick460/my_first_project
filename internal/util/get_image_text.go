package util

import (
	"github.com/otiai10/gosseract/v2"
)

func GetImageText(imagePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage(imagePath)
	if err != nil {
		return "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}
