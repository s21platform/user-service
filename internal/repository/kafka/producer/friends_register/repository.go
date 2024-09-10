package friends_register

import (
	"context"
	"github.com/s21platform/user-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type FriendsInvite struct {
	pdsr *kafka.Writer
}

func New(cfg *config.Config) *FriendsInvite {

	pdsr := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Broker),
		Topic:        cfg.Kafka.FriendsRegister,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &FriendsInvite{pdsr: pdsr}
}

func (f FriendsInvite) SendMessage(ctx context.Context, email string) error {
	err := f.pdsr.WriteMessages(ctx, kafka.Message{
		Value: []byte(email),
	})
	if err != nil {
		return err
	}
	return nil
}
