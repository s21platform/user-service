package new_avatar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/s21platform/user-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	consumer *kafka.Reader
	dbR      DBRepo
}

type AvatarUpdateRsvMap struct {
	UUID string `json:"uuid"`
	Link string `json:"link"`
}

func New(cfg *config.Config, db DBRepo) (*KafkaConsumer, error) {
	if cfg.Kafka.Server == "" {
		return nil, fmt.Errorf("error cfg.Kafka.Server is empty")
	}

	if cfg.Kafka.SetNewAvatar == "" {
		return nil, fmt.Errorf("error cfg.Kafka.SetNewAvatar is empty")
	}

	broker := []string{cfg.Kafka.Server}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   cfg.Kafka.SetNewAvatar,
		GroupID: "123",
	})

	return &KafkaConsumer{consumer: reader, dbR: db}, nil
}

func (kc *KafkaConsumer) Listen() {
	for { // TODO подумать над запуском в go рутине
		readMsg, err := kc.readMessage()

		if err != nil {
			log.Println("error kc.readMessage: ", err)
			continue
		}

		err = kc.dbR.UpdateUserAvatar(readMsg.UUID, readMsg.Link)

		if err != nil {
			log.Println("error kc.dbR.UpdateUserAvatar: ", err)
			continue
		}
	}
}

func (kc *KafkaConsumer) readMessage() (AvatarUpdateRsvMap, error) {
	var Avatar AvatarUpdateRsvMap

	msgJSON, err := kc.consumer.ReadMessage(context.Background())

	if err != nil {
		return Avatar, err
	}

	err = json.Unmarshal(msgJSON.Value, &Avatar)

	if err != nil {
		return Avatar, err
	}

	log.Println("read from topic (avatarLink):", Avatar.Link)

	return Avatar, nil
}
