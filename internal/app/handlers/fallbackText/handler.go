package fallbackText

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"notifier-bot-telegram/internal/internalLogger"
	"notifier-bot-telegram/pkg/logger"
)

type telegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Handler struct {
	telegramClient telegramClient
	log            internalLogger.Logger
}

func New(telegramClient telegramClient, log internalLogger.Logger) *Handler {
	return &Handler{
		telegramClient: telegramClient,
		log:            log,
	}
}

func (h *Handler) Handle(ctx context.Context, update tgbotapi.Update) error {
	// TODO вынести обогащение служебными полями в consumer
	logCtx := logger.WithAttrs(ctx, map[string]any{
		"updateID":  update.UpdateID,
		"messageID": update.Message.MessageID,
		"command":   update.Message.Command(),
		"from":      update.Message.From,
		"chatID":    update.Message.Chat.ID,
	})
	if err := h.telegramClient.SendMessage(logCtx, update.Message.Chat.ID, "В данный момент я не способен поддержать диалог.\n Воспользуйтесь /help, чтобы посмотреть список доступных команд. "); err != nil {
		h.log.ErrorContext(logCtx, "error while handle fallback command", err)
		return err
	}
	h.log.InfoContext(logCtx, "fallback command successfully handled")
	return nil
}
