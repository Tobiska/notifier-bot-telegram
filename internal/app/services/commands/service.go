package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	"notifier-bot-telegram/internal/app/models"
	detailsRepository "notifier-bot-telegram/internal/app/repository/details"
)

type userRepository interface {
	SaveUser(ctx context.Context, user models.User) error
	FindUser(ctx context.Context, chatID int64) (*models.User, error)
	UpdateStatus(ctx context.Context, chatID int64, status models.Status) error
}

type detailRepository interface {
	Begin() (*sql.Tx, error)
	WithTx(tx *sql.Tx) *detailsRepository.Repository
	PartialUpdate(ctx context.Context, detailID int64, updateDetail detailsRepository.UpdateModel) error
	Save(ctx context.Context, detail models.Detail) error
}

type telegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Service struct {
	userRepository   userRepository
	detailRepository detailRepository
	telegramClient   telegramClient
}

func New(userRepository userRepository, detailRepository detailRepository, telegramClient telegramClient) *Service {
	return &Service{
		userRepository:   userRepository,
		telegramClient:   telegramClient,
		detailRepository: detailRepository,
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

func (s *Service) ApplyStartAdd(ctx context.Context, chatID int64) error {
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

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.AddName); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}

		return nil
	}

	return nil
}

func (s *Service) Add(ctx context.Context, chatID int64, text string) error {
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

	if u.Status == models.AddName {
		if err := s.saveOrUpdateDetails(ctx, chatID, detailsRepository.UpdateModel{
			Name: pointer.To(text),
		}); err != nil {
			return fmt.Errorf("save details error: %w", err)
		}

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.AddSoftDeadline); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}

		if err := s.telegramClient.SendMessage(ctx, chatID, "1. Напишите первого предупреждения о тех. обслуживании в формате `2006-01-02`."); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
		return nil
	}

	if u.Status == models.AddSoftDeadline {
		timestamp, err := time.Parse(time.DateOnly, text)
		if err != nil {
			if err := s.telegramClient.SendMessage(ctx, chatID, fmt.Sprintf("Кажется вы прислали дату в некоректном формате.")); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return fmt.Errorf("can't parse text for softdeadline: %w", err)
		}

		if err := s.saveOrUpdateDetails(ctx, chatID, detailsRepository.UpdateModel{
			SoftDeadline: pointer.To(timestamp),
		}); err != nil {
			return fmt.Errorf("save details error: %w", err)
		}

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.AddHardDeadline); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}

		if err := s.telegramClient.SendMessage(ctx, chatID, "1. Напишите второго предупреждения о тех. обслуживании в формате `2006-01-02`."); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	if u.Status == models.AddHardDeadline {
		timestamp, err := time.Parse(time.DateOnly, text)
		if err != nil {
			if err := s.telegramClient.SendMessage(ctx, chatID, fmt.Sprintf("Кажется вы прислали дату в некоректном формате.")); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return fmt.Errorf("can't parse text for harddeadline: %w", err)
		}

		if err := s.saveOrUpdateDetails(ctx, chatID, detailsRepository.UpdateModel{
			HardDeadline: pointer.To(timestamp),
		}); err != nil {
			return fmt.Errorf("save details error: %w", err)
		}

		if err := s.userRepository.UpdateStatus(ctx, chatID, models.Wait); err != nil {
			return fmt.Errorf("error while update message: %w", err)
		}

		if err := s.telegramClient.SendMessage(ctx, chatID, "Всё готово! Я напомню вам о тех обслуживании, когда наступит время."); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	return nil
}

func (s *Service) saveOrUpdateDetails(ctx context.Context, chatID int64, updateModel detailsRepository.UpdateModel) error {
	var err error
	tx, err := s.detailRepository.Begin()
	if err != nil {
		return fmt.Errorf("begin tx error: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txDetailsRepo := s.detailRepository.WithTx(tx)

	details, err := txDetailsRepo.Search(ctx, detailsRepository.Filter{
		ChatID:   pointer.ToInt64(chatID),
		IsFilled: pointer.ToBool(true),
	})
	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}
	if len(details) == 0 {
		err = txDetailsRepo.Save(ctx, models.Detail{ChatID: chatID, Name: *updateModel.Name})
		if err != nil {
			return fmt.Errorf("save details error: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}

	detail := details[0]

	err = txDetailsRepo.PartialUpdate(ctx, detail.ID, updateModel)
	if err != nil {
		return fmt.Errorf("partial update details error: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
