package messages

import (
	"context"
	"demochat/clients/insights"
	"demochat/internal/repositories/messages"
	"log"
)

type Service interface {
	SaveMessage(ctx context.Context, userID int64, message string) error
	ListMessages(ctx context.Context) ([]*messages.Message, error)
}

type service struct {
	repo           messages.Repository
	insightsClient insights.Client
}

func New(repo messages.Repository, ic insights.Client) Service {
	return &service{repo: repo, insightsClient: ic}
}

func (s *service) SaveMessage(ctx context.Context, userID int64, message string) error {
	msg := &messages.Message{
		UserID:  userID,
		Message: message,
	}
	err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return err
	}

	data := insights.PublishMessageRequest{UserID: userID, Message: message}

	go func() {
		log.Println("attempting to publish message")
		err := s.insightsClient.PublishMessage(ctx, data)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (s *service) ListMessages(ctx context.Context) ([]*messages.Message, error) {
	return s.repo.ListMessages(ctx)
}
