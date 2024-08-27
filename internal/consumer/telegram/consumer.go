package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type router interface {
	Route(ctx context.Context, messages []tgbotapi.Message) error
}

type updater interface {
	GetUpdates() tgbotapi.UpdatesChannel
}

type Consumer struct {
	router  router
	updater updater
}

func New(router router, updater updater) *Consumer {
	return &Consumer{
		router:  router,
		updater: updater,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for update := range c.updater.GetUpdates() {
		if update.Message != nil {

		}
	}
	return nil
}

func (c *Consumer) handle(ctx context.Context, message tgbotapi.Message) error {
	return nil
}
