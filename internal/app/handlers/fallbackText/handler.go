package fallbackText

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"notifier-bot-telegram/internal/internalLogger"
	"notifier-bot-telegram/pkg/logger"
)

type telegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type service interface {
	Add(ctx context.Context, chatID int64, text string) error
}

type Handler struct {
	telegramClient telegramClient
	service        service
	log            internalLogger.Logger
}

func New(telegramClient telegramClient, service service, log internalLogger.Logger) *Handler {
	return &Handler{
		telegramClient: telegramClient,
		service:        service,
		log:            log,
	}
}

func (h *Handler) Handle(ctx context.Context, update tgbotapi.Update) error {
	if err := validate(update); err != nil {
		h.log.WarnContext(ctx, "validate error", err)
		return err
	}

	// TODO вынести обогащение служебными полями в consumer
	logCtx := logger.WithAttrs(ctx, map[string]any{
		"updateID":  update.UpdateID,
		"messageID": update.Message.MessageID,
		"command":   update.Message.Command(),
		"from":      update.Message.From,
		"chatID":    update.Message.Chat.ID,
	})

	if err := h.service.Add(ctx, update.Message.Chat.ID, update.Message.Text); err != nil {
		h.log.ErrorContext(logCtx, "fallback command successfully handled", err)
		return fmt.Errorf("error handle update: %w", err)
	}

	h.log.InfoContext(logCtx, "fallback command successfully handled")
	return nil
}

func validate(update tgbotapi.Update) error {
	if update.Message == nil {
		return errors.New("message is empty")
	}

	if update.Message.Chat == nil {
		return errors.New("chat is empty")
	}

	return nil
}
