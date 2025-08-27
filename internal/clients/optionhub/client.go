package optoinhub

import (
	"context"
	"fmt"
	"log"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	optionhub "github.com/s21platform/optionhub-service/pkg/optionhub"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
)

type Client struct {
	client optionhub.OptionhubServiceClient
}

func (c *Client) GetAttributesMeta(ctx context.Context, attributeIds []model.Attribute) (model.AttributeMetaMap, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	out, err := c.client.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
		EntityAttributeIds: lo.Map(attributeIds, func(item model.Attribute, _ int) int64 {
			return item.Int64()
		}),
	})
	if err != nil {
		return nil, err
	}
	res := make(model.AttributeMetaMap)
	for key, value := range out.AttributesMetadata {
		res[key] = model.AttributeMeta{
			Label: value.Label,
			Type:  value.Type.String(),
		}
	}
	return res, nil
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Optionhub.Host, cfg.Optionhub.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to community service: %v", err)
	}
	client := optionhub.NewOptionhubServiceClient(conn)
	return &Client{client: client}
}
