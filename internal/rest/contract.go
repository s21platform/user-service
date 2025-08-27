package rest

import (
	"context"

	"github.com/s21platform/user-service/internal/model"
)

type DbRepo interface {
	GetNicknameByUuid(ctx context.Context, uuid string) (string, error)
	GetPersonalityByUuid(ctx context.Context, uuid string) (model.Personality, error)
}

type OptionhubClient interface {
	GetAttributesMeta(ctx context.Context, attributeIds []model.Attribute) (model.AttributeMetaMap, error)
}
