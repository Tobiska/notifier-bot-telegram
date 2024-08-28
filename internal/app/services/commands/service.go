package commands

import (
	"context"
	"fmt"
	"notifier-bot-telegram/internal/app/models"
	"time"
)

type userRepository interface {
	SaveUser(ctx context.Context, user models.User) error
	FindUser(ctx context.Context, chatID int64) (*models.User, error)
	UpdateStatus(ctx context.Context, chatID int64, status models.Status) error
}

type telegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Service struct {
	userRepository userRepository
	telegramClient telegramClient
}

func New(userRepository userRepository, telegramClient telegramClient) *Service {
	return &Service{
		userRepository: userRepository,
		telegramClient: telegramClient,
	}
}

func (s *Service) Start(ctx context.Context, user models.User) error {
	existUser, err := s.userRepository.FindUser(ctx, user.ChatID)
	if err != nil {
		return fmt.Errorf("error while find user: %w", err)
	}

	if existUser != nil {
		if err := s.telegramClient.SendMessage(ctx, user.ChatID, fmt.Sprintf("Похоже раньше вы уже пользовались бото. Аккаунт был зарегестрирован %s", existUser.CreatedAt.Format(time.DateOnly))); err != nil {
			return fmt.Errorf("error while send error already started: %w", err)
		}
		return nil
	}

	user.Status = models.Created
	if err := s.userRepository.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("error while save user :%w", err)
	}

	if err := s.telegramClient.SendMessage(ctx, user.ChatID, "Привет! Я TruckFixBot я умею напоминать о своевременном техническом обслуживании.\n Для того, чтобы новую добавить деталь для отслеживания используй команду /add."); err != nil {
		return fmt.Errorf("error while send error already started: %w", err)
	}

	return nil
}

func (s *Service) Add(ctx context.Context, chatID int64) error {
	u, err := s.userRepository.FindUser(ctx, chatID)
	if err != nil {
		return fmt.Errorf("error while find user: %w", err)
	}

	if u == nil {
		if err := s.telegramClient.SendMessage(ctx, chatID, "Кажется вы ещё не пользовались ботов, попробуйте выполнить команду /start."); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
		return nil
	}

	if u.Status == models.Created || u.Status == models.Wait {
		if err := s.telegramClient.SendMessage(ctx, chatID, "1. Напишите имя детали"); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.StartAdded); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}
		return nil
	}

	if u.Status == models.StartAdded {
		if err := s.telegramClient.SendMessage(ctx, chatID, "1. Напишите имя детали"); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.StartAdded); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}
		return nil
	}

	return nil
}
