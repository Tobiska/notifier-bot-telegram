package cmd

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"notifier-bot-telegram/internal/clients/telegram"
	"notifier-bot-telegram/internal/config"
)

func run() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	telegramClient, err := telegram.New(cfg.Telegram)
	if err != nil {
		return fmt.Errorf("create telegram client error: %w", err)
	}

	for update := range telegramClient.GetUpdates() {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			if err := telegramClient.SendMessage(context.Background(), update.Message.Chat.ID, update.Message.Text); err != nil {
				return fmt.Errorf("send message error: %w", err)
			}
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
