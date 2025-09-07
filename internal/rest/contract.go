package rest

import (
	"context"

	"github.com/s21platform/user-service/internal/model"
)

type DbRepo interface {
	GetNicknameByUuid(ctx context.Context, uuid string) (string, error)
	GetPersonalityByUuid(ctx context.Context, uuid string) (model.Personality, error)
	GetUserAttributesByUuid(ctx context.Context, uuid string) (model.UserAttributes, error)
	UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error
}

type OptionhubClient interface {
	GetAttributesMeta(ctx context.Context, attributeIds []model.Attribute) (model.AttributeMetaMap, error)
}
