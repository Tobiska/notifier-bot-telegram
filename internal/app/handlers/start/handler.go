package start

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"notifier-bot-telegram/internal/app/models"
	"notifier-bot-telegram/internal/internalLogger"
	"notifier-bot-telegram/pkg/logger"
)

type service interface {
	Start(ctx context.Context, user models.User) error
}

type Handler struct {
	service service
	log     internalLogger.Logger
}

func New(service service, log internalLogger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
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
	if err := h.service.Start(logCtx, models.User{
		ChatID:   update.Message.Chat.ID,
		UserID:   update.Message.From.ID,
		Username: update.Message.From.UserName,
	}); err != nil {
		h.log.ErrorContext(logCtx, "error while handle fallback command", err)
		return err
	}
	h.log.InfoContext(logCtx, "fallback command successfully handled")
	return nil
}
