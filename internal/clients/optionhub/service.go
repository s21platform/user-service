package optoinhub

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
)

type Handle struct {
	client optionhubproto.OptionhubServiceClient
}

func (h *Handle) GetOs(ctx context.Context, id *int64) (*model.OS, error) {
	if id == nil {
		return nil, nil
	}

	os, err := h.client.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: *id})
	if err != nil {
		return nil, err
	}
	return &model.OS{
		Id:    os.Id,
		Label: os.Value,
	}, nil
}

func MustConnect(cfg *config.Config) *Handle {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Optionhub.Host, cfg.Optionhub.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to community service: %v", err)
	}
	client := optionhubproto.NewOptionhubServiceClient(conn)
	return &Handle{client: client}
}
