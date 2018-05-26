package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	msg.ParseMode = "Markdown"
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
// ALARM! GOVNOCODE!
// TODO: refactor
func FormatMessage(dceInfo []string, diff string) string {
	var name = dceInfo[0]
	var dceLink, tradeLink, pairsInfo, message string
	pairsInfo = diff
	if len(dceInfo) == 2 {
		dceLink = dceInfo[1]
	}
	if len(dceInfo) == 3 {
		pairsInfo = ""
		dceLink = dceInfo[1]
		tradeLink = dceInfo[2]
		diff = strings.Trim(diff, "\n")
		diffs := strings.Split(diff, "\n")
		for _, pair := range diffs {
			currency := strings.Split(trimLeftChars(pair, 1), "-")
			pairsInfo += fmt.Sprintf("[%v]("+tradeLink+")\n", pair, strings.Trim(currency[0], " "), currency[1])
		}
	}
	message = fmt.Sprintf("[%v](%v)* has pairs alerts:*\n\n%v", name, dceLink, pairsInfo)
	return message
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}
