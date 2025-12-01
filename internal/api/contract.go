package api

import (
	"context"

	"github.com/google/uuid"

	"github.com/s21platform/user-service/internal/model"
)

type DbRepo interface {
	GetNicknameByUuid(ctx context.Context, uuid string) (string, error)
	GetPersonalityByUuid(ctx context.Context, uuid string) (model.Personality, error)
	GetUserAttributesByUuid(ctx context.Context, uuid string) (model.UserAttributes, error)
	UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error
	CreatePost(ctx context.Context, uuid uuid.UUID, content string) (string, error)
	EditPost(ctx context.Context, postID string, ownerUUID uuid.UUID, content string) (string, error)
	DeletePost(ctx context.Context, postID string, ownerUUID uuid.UUID) error
	CheckPostOwnership(ctx context.Context, postID string, ownerUUID uuid.UUID) error
}

type OptionhubClient interface {
	GetAttributesMeta(ctx context.Context, attributeIds []model.Attribute) (model.AttributeMetaMap, error)
}

type PostCreatedProducer interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}
