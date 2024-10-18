package optoinhub

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"log"

	"github.com/s21platform/user-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

type Handle struct {
	client optionhubproto.OptionhubServiceClient
}

func (h *Handle) GetOs(ctx context.Context, id *int64) (*string, error) {
	if id == nil {
		return nil, nil
	}

	os, err := h.client.GetOsById(ctx, &optionhubproto.GetByIdIn{Id: *id})
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(os.Value), nil
}

func MustConnect(cfg *config.Config) *Handle {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Optionhub.Host, cfg.Optionhub.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to community service: %v", err)
	}
	Client := optionhubproto.NewOptionhubServiceClient(conn)
	return &Handle{client: Client}
}
