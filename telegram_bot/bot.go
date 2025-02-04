package telegram_bot

import (
	"fmt"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Send(bot *tgbotapi.BotAPI, id int64, message string, messageType string) error {
	var msg tgbotapi.Chattable

	switch messageType {
	case "text":
		msg = tgbotapi.NewMessage(id, message)
	case "animation":
		msg = tgbotapi.NewAnimationShare(id, message)
	case "sticker":
		msg = tgbotapi.NewStickerShare(id, message)
	default:
		return fmt.Errorf("unsupported message type: %v", messageType)
	}

	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

func SendTelegramMsg(token string, id int64, chall string, user string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("error creating bot: %v", err)
	}

	message := messages[rand.Intn(len(messages))]

	messageText := strings.ReplaceAll(message.Text, "<user>", user)
	messageText = strings.ReplaceAll(messageText, "<chall>", chall)

	err = Send(bot, id, messageText, "text")
	if err != nil {
		return err
	}

	err = Send(bot, id, message.Media, message.MediaType)
	if err != nil {
		return err
	}

	return nil
}
