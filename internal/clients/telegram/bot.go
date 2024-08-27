package telegram

import "C"
import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"notifier-bot-telegram/internal/config"
)

type Client struct {
	botApi *tgbotapi.BotAPI
	cfg    config.Telegram
}

func New(cfg config.Telegram) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotFatherToken)
	if err != nil {
		return nil, fmt.Errorf("error while creating bot api: %w", err)
	}
	bot.Debug = cfg.Debug

	return &Client{
		botApi: bot,
	}, nil
}

func (c *Client) GetUpdates() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0) // update offset
	updateConfig.Timeout = c.cfg.TimeoutUpdates

	return c.botApi.GetUpdatesChan(updateConfig)
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.botApi.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
