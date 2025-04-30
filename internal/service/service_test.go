package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/distribution/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
	"github.com/s21platform/user-service/pkg/user"
)

func TestServer_GetUsersByUUID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockDBRepo := NewMockDbRepo(ctrl)

	t.Run("should ok", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate().String()
		mockInput := &user.GetUsersByUUIDIn{
			UsersUuid: []*user.UsersUUID{
				{Uuid: peer1},
				{Uuid: peer2},
			},
		}

		mockDBRepo.EXPECT().GetUsersByUUID(peer1).Return(model.UserInfoMin{
			Uuid:       peer1,
			Login:      "login1",
			LastAvatar: "avatar1",
			Name:       "John",
			Surname:    "Doe",
		}, nil)

		mockDBRepo.EXPECT().GetUsersByUUID(peer2).Return(model.UserInfoMin{
			Uuid:       peer2,
			Login:      "login2",
			LastAvatar: "avatar2",
			Name:       "Jane",
			Surname:    "Doe",
		}, nil)
		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetUsersByUUID(ctx, mockInput)
		assert.NoError(t, err)
		assert.Equal(t, result, &user.GetUsersByUUIDOut{
			UsersInfo: []*user.UserInfoMin{
				{
					Uuid:       peer1,
					Login:      "login1",
					LastAvatar: "avatar1",
					Name:       "John",
					Surname:    "Doe",
				},
				{
					Uuid:       peer2,
					Login:      "login2",
					LastAvatar: "avatar2",
					Name:       "Jane",
					Surname:    "Doe",
				},
			},
		})
	})
	t.Run("should_return_nil_if_empty_UUID_provided", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		mockInput := &user.GetUsersByUUIDIn{}

		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetUsersByUUID(ctx, mockInput)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("should_return_if_uuid_is_empty", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate().String()
		mockInput := &user.GetUsersByUUIDIn{
			UsersUuid: []*user.UsersUUID{
				{Uuid: peer1},
				{},
				{Uuid: peer2},
			},
		}

		mockDBRepo.EXPECT().GetUsersByUUID(peer1).Return(model.UserInfoMin{
			Uuid:       peer1,
			Login:      "login1",
			LastAvatar: "avatar1",
			Name:       "John",
			Surname:    "Doe",
		}, nil)

		mockDBRepo.EXPECT().GetUsersByUUID(peer2).Return(model.UserInfoMin{
			Uuid:       peer2,
			Login:      "login2",
			LastAvatar: "avatar2",
			Name:       "Jane",
			Surname:    "Doe",
		}, nil)

		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetUsersByUUID(ctx, mockInput)
		assert.NoError(t, err)
		assert.Equal(t, result, &user.GetUsersByUUIDOut{
			UsersInfo: []*user.UserInfoMin{
				{
					Uuid:       peer1,
					Login:      "login1",
					LastAvatar: "avatar1",
					Name:       "John",
					Surname:    "Doe",
				},
			},
		})
	})
	t.Run("should_return_error_if_db_fails", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)

		mockInput := &user.GetUsersByUUIDIn{
			UsersUuid: []*user.UsersUUID{
				{Uuid: peer1},
			},
		}

		expectedError := fmt.Errorf("database error")
		mockDBRepo.EXPECT().GetUsersByUUID(peer1).Return(model.UserInfoMin{}, expectedError)

		s := &Server{dbRepo: mockDBRepo}

		result, err := s.GetUsersByUUID(ctx, mockInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user by UUID")
		assert.Contains(t, err.Error(), peer1)
		assert.Contains(t, err.Error(), expectedError.Error())
		assert.Nil(t, result)
	})
}
