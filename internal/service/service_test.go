package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
		peer1 := uuid.New().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.New().String()
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
		peer1 := uuid.New().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		mockInput := &user.GetUsersByUUIDIn{}

		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetUsersByUUID(ctx, mockInput)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("should_return_if_uuid_is_empty", func(t *testing.T) {
		peer1 := uuid.New().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.New().String()
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
		peer1 := uuid.New().String()
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

func TestServer_SetFriends(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(false, nil)
		mockDBRepo.EXPECT().SetFriends(ctx, peer1, peer2).Return(nil)
		mockLogger.EXPECT().AddFuncName("SetFriends")

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})
		assert.NoError(t, err)
		assert.Equal(t, res, &user.SetFriendsOut{Success: true})
	})
	t.Run("should_is_friends", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(true, nil)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error("user already in friends")

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})

		assert.NoError(t, err)
		assert.Equal(t, &user.SetFriendsOut{Success: true}, res)
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error("failed to get user UUID from context")

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: uuid.New().String()})

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get user UUID from context")
	})

	t.Run("should_repo_err", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()
		repoErr := errors.New("test error")

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(false, nil)
		mockDBRepo.EXPECT().SetFriends(ctx, peer1, peer2).Return(repoErr)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error("failed to SetFriends from dbRepo")

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
	})

	t.Run("should_check_friendship_error", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()
		checkErr := errors.New("check error")

		ctx := context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(false, checkErr)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error(" failed to check user friendship")

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user friendship")
	})
}
