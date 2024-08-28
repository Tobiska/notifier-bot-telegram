package telegram

import (
	"context"
	rout "notifier-bot-telegram/internal/router"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type router interface {
	Route(ctx context.Context, update tgbotapi.Update) ([]rout.Handler, error)
}

type updater interface {
	GetUpdates() tgbotapi.UpdatesChannel
}

type retry interface {
	DoRetry(func() error) error
}

type Consumer struct {
	router  router
	updater updater
	retry   retry
}

func New(router router, updater updater, retry retry) *Consumer {
	return &Consumer{
		router:  router,
		updater: updater,
		retry:   retry,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for update := range c.updater.GetUpdates() {
		if err := c.handleUpdate(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (c *Consumer) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	handlers, err := c.router.Route(ctx, update)
	if err != nil {
		return err
	}

	for _, h := range handlers {
		err := c.retry.DoRetry(func() error {
			return h.Handle(ctx, update)
		})
		if err != nil {
			return err
		}
	}

	return err
}
