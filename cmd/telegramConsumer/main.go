package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"notifier-bot-telegram/internal/app/handlers/add"
	"notifier-bot-telegram/internal/app/handlers/fallbackCommand"
	"notifier-bot-telegram/internal/app/handlers/fallbackText"
	"notifier-bot-telegram/internal/app/handlers/start"
	"notifier-bot-telegram/internal/app/repository/details"
	"notifier-bot-telegram/internal/app/repository/user"
	"notifier-bot-telegram/internal/app/services/commands"
	telegramConsumer "notifier-bot-telegram/internal/consumer/telegram"
	"notifier-bot-telegram/internal/router"
	"notifier-bot-telegram/pkg/logger"
	"notifier-bot-telegram/pkg/sqlite"
	sql2 "notifier-bot-telegram/pkg/storage/sql"
	"notifier-bot-telegram/pkg/utils"

	"notifier-bot-telegram/internal/clients/telegram"
	"notifier-bot-telegram/internal/config"
)

func run() error {
	ctx := context.Background()

	log := logger.NewLogger(logger.WithJSONHandler())

	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	retry := utils.NewRetry(cfg.Telegram.RetryCount, cfg.Telegram.RetryTimeout)

	telegramClient, err := telegram.New(cfg.Telegram)
	if err != nil {
		return fmt.Errorf("create telegram client error: %w", err)
	}

	db, err := sqlite.New(cfg.Dsn)
	if err != nil {
		return fmt.Errorf("failed connect to database: %w", err)
	}
	defer db.Close()

	sqlStorage := sql2.New(db)

	userRepository := user.New(sqlStorage)
	detailsRepository := details.New(sqlStorage)

	commandsService := commands.New(userRepository, detailsRepository, telegramClient)

	fallbackCommandHandler := fallbackCommand.New(telegramClient, log)
	fallbackTextHandler := fallbackText.New(telegramClient, commandsService, log)

	handlerRouter := router.New(fallbackCommandHandler, fallbackTextHandler)
	handlerRouter.RegisterCommand("/start", start.New(commandsService, log))
	handlerRouter.RegisterCommand("/add", add.New(commandsService, log))

	telegramConsumer := telegramConsumer.New(handlerRouter, telegramClient, retry)

	notifyCtx, cancel := signal.NotifyContext(ctx, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	logCtx := logger.CtxWithSystemAttrs(notifyCtx)

	g, errGroupCtx := errgroup.WithContext(logCtx)

	g.Go(func() error {
		if err := telegramConsumer.Run(errGroupCtx); err != nil {
			log.ErrorContext(ctx, "fatal error telegram consumer", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
