package bot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/telegram-bot-api.v4"
)

// TelegramConfig is a struct with config data.
// Initialized by GetTelegramConfig func
type TelegramConfig struct {
	Token  string
	ChatID int64
}

// SendMessageToTelegramChannel do what it has in its name )
// Need to provide token and chatID for a group
// Here is a method how to get chatID for a private group:
// https://stackoverflow.com/questions/33858927/how-to-obtain-the-chat-id-of-a-private-telegram-channel
func SendMessageToTelegramChannel(config TelegramConfig, message string) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(config.ChatID, message)
	bot.Send(msg)
	return nil
}

// GetTelegramConfig returns TelegramConfig struct
// filled with data from .env file
func GetTelegramConfig(envFilePath string) (TelegramConfig, error) {
	var err error
	if envFilePath != "" {
		err = godotenv.Load(envFilePath)
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		return TelegramConfig{}, err
	}
	token := os.Getenv("TOKEN")
	chatID, err := strconv.ParseInt(os.Getenv("CHATID"), 10, 64)
	if err != nil {
		return TelegramConfig{}, err
	}
	return TelegramConfig{Token: token, ChatID: chatID}, nil
}

// FormatMessage takes two arguments and for a message which
// the bot can send to a channel
func FormatMessage(dceInfo []string, diff string) string {
	var name = dceInfo[0]
	var link, message string
	if len(dceInfo) == 2 {
		link = dceInfo[1]
	}
	message = fmt.Sprintf("%v %v\n%v", name, link, diff)
	return message
}
