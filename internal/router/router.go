package router

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, update tgbotapi.Update) error
}

type Router struct {
	commandHandlers        map[string]Handler
	fallbackCommandHandler Handler
	textHandlers           map[string]Handler
	fallbackTextHandler    Handler
	mu                     *sync.RWMutex
}

func New(fallbackCommandHandler, fallbackTextHandler Handler) *Router {
	return &Router{
		commandHandlers:        make(map[string]Handler),
		fallbackCommandHandler: fallbackCommandHandler,
		fallbackTextHandler:    fallbackTextHandler,
		mu:                     &sync.RWMutex{},
	}
}

func (r *Router) RegisterCommand(path string, handler Handler) {
	if path == "" {
		panic("path must be not empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.commandHandlers[path]; ok {
		panic(fmt.Sprintf("handler with path: %s already exist", path))
	}

	r.commandHandlers[path] = handler
}

func (r *Router) Route(ctx context.Context, update tgbotapi.Update) ([]Handler, error) {
	var handlers []Handler
	if update.Message != nil {
		messageHandler, err := r.routeMessage(ctx, *update.Message)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, messageHandler)
	}

	return handlers, nil
}

func (r *Router) routeMessage(_ context.Context, msg tgbotapi.Message) (Handler, error) {
	if msg.IsCommand() {
		r.mu.RLock()
		defer r.mu.RUnlock()
		if handler, ok := r.commandHandlers[msg.Text]; ok {
			return handler, nil
		}
		return r.fallbackCommandHandler, nil
	} else {
		r.mu.RLock()
		defer r.mu.RUnlock()
		if handler, ok := r.textHandlers[msg.Text]; ok {
			return handler, nil
		}
		return r.fallbackTextHandler, nil
	}
}
