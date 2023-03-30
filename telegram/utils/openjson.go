package utils

import (
	"encoding/json"
	"github.com/MamushevArup/telegram-bot-krisha/internal/texts"
	"os"
)

func OpenJsonFile(filepath string) (string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return "Cannot open json file", err
	}
	var config texts.Config
	defer file.Close()
	decode := json.NewDecoder(file)
	err = decode.Decode(&config)
	res := config.Token
	return res, nil
}
