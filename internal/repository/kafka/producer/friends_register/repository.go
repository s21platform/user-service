package friends_register

import (
	"context"
	"encoding/json"
	"github.com/s21platform/user-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type FriendsInvite struct {
	pdsr *kafka.Writer
}

type FriendsInviteSendMap struct {
	Email string `json:"email"`
	UUID  string `json:"uuid"`
}

func New(cfg *config.Config) *FriendsInvite {

	pdsr := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Server),
		Topic:        cfg.Kafka.FriendsRegister,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &FriendsInvite{pdsr: pdsr}
}

func (f FriendsInvite) SendMessage(ctx context.Context, email string, uuid string) error {
	mess := FriendsInviteSendMap{
		Email: email,
		UUID:  uuid,
	}

	messJson, err := json.Marshal(mess)
	if err != nil {
		return err
	}

	err = f.pdsr.WriteMessages(ctx, kafka.Message{
		Value: messJson,
	})
	if err != nil {
		return err
	}
	return nil
}
