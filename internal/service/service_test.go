package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/guregu/null/v6"

	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})

		assert.NoError(t, err)
		assert.Equal(t, &user.SetFriendsOut{Success: true}, res)
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

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

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.SetFriends(ctx, &user.SetFriendsIn{Peer: peer2})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user friendship")
	})
}

func TestServer_RemoveFriends(t *testing.T) {
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
		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(true, nil)
		mockDBRepo.EXPECT().RemoveFriends(ctx, peer1, peer2).Return(nil)
		s := &Server{dbRepo: mockDBRepo}
		_, err := s.RemoveFriends(ctx, &user.RemoveFriendsIn{Peer: peer2})
		assert.NoError(t, err)
	})
	t.Run("should_no_UUID", func(t *testing.T) {
		ctx = context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		s := &Server{dbRepo: mockDBRepo}

		_, err := s.RemoveFriends(ctx, &user.RemoveFriendsIn{Peer: uuid.New().String()})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "failed to get user UUID from context")
	})

	t.Run("should_error_when_not_friends", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(false, nil)

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.RemoveFriends(ctx, &user.RemoveFriendsIn{Peer: peer2})

		assert.NoError(t, err)
		assert.False(t, res.Success)
	})
	t.Run("should_error_when_check_friendship_fails", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()
		expectedErr := errors.New("check friendship error")

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(false, expectedErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.RemoveFriends(ctx, &user.RemoveFriendsIn{Peer: peer2})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user friendshipt")
	})

	t.Run("should_error_when_remove_friends_fails", func(t *testing.T) {
		peer1 := uuid.New().String()
		peer2 := uuid.New().String()
		expectedErr := errors.New("remove friends error")

		ctx = context.WithValue(context.Background(), config.KeyUUID, peer1)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, peer1, peer2).Return(true, nil)
		mockDBRepo.EXPECT().RemoveFriends(ctx, peer1, peer2).Return(expectedErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.RemoveFriends(ctx, &user.RemoveFriendsIn{Peer: peer2})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestServer_GetCountFriends(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	t.Run("should_ok", func(t *testing.T) {
		peer := uuid.New().String()
		var subscription int64 = 10
		var subscribers int64 = 5
		ctx = context.WithValue(context.Background(), config.KeyUUID, peer)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().GetSubscriptionCount(ctx, peer).Return(subscription, nil)
		mockDBRepo.EXPECT().GetSubscribersCount(ctx, peer).Return(subscribers, nil)

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.GetCountFriends(ctx, &user.EmptyFriends{})

		assert.NoError(t, err)
		assert.Equal(t, &user.GetCountFriendsOut{
			Subscription: subscription,
			Subscribers:  subscribers,
		}, res)
	})
	t.Run("should_error_when_no_UUID_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetCountFriends(ctx, &user.EmptyFriends{})

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get user UUID from context")
	})
	t.Run("should_error_when_get_subscription_fails", func(t *testing.T) {
		peer := uuid.New().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, peer)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		expectedErr := errors.New("subscription error")

		mockDBRepo.EXPECT().GetSubscriptionCount(ctx, peer).Return(int64(0), expectedErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetCountFriends(ctx, &user.EmptyFriends{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("should_error_when_get_subscribers_fails", func(t *testing.T) {
		peer := uuid.New().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, peer)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		expectedErr := errors.New("subscribers error")

		mockDBRepo.EXPECT().GetSubscriptionCount(ctx, peer).Return(int64(10), nil)
		mockDBRepo.EXPECT().GetSubscribersCount(ctx, peer).Return(int64(0), expectedErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetCountFriends(ctx, &user.EmptyFriends{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestServer_GetPeerFollow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	t.Run("should_ok", func(t *testing.T) {
		userUUID := uuid.New().String()
		ctx = context.WithValue(context.Background(), config.KeyUUID, userUUID)
		followersUUID := []string{
			uuid.New().String(),
			uuid.New().String(),
		}
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().GetPeerFollow(ctx, userUUID).Return(followersUUID, nil)

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.GetPeerFollow(ctx, &user.GetPeerFollowIn{Uuid: userUUID})
		assert.NoError(t, err)
		assert.Equal(t, res, &user.GetPeerFollowOut{Subscription: []*user.Peer{
			{Uuid: followersUUID[0]},
			{Uuid: followersUUID[1]},
		}})
	})
	t.Run("should_error_when_no_UUID_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetPeerFollow(ctx, &user.GetPeerFollowIn{
			Uuid: uuid.New().String(),
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get user UUID from context")
	})
	t.Run("should_repo_err", func(t *testing.T) {
		userUUID := uuid.New().String()
		repoErr := errors.New("test")

		ctx = context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().GetPeerFollow(ctx, userUUID).Return(nil, repoErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetPeerFollow(ctx, &user.GetPeerFollowIn{Uuid: userUUID})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_GetWhoFollowPeer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	t.Run("should_ok", func(t *testing.T) {
		userUUID := uuid.New().String()
		ctx = context.WithValue(context.Background(), config.KeyUUID, userUUID)
		followersUUID := []string{
			uuid.New().String(),
			uuid.New().String(),
		}
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().GetWhoFollowPeer(ctx, userUUID).Return(followersUUID, nil)
		s := &Server{dbRepo: mockDBRepo}
		res, err := s.GetWhoFollowPeer(ctx, &user.GetWhoFollowPeerIn{
			Uuid: userUUID,
		})
		assert.NoError(t, err)
		assert.Equal(t, res, &user.GetWhoFollowPeerOut{
			Subscribers: []*user.Peer{
				{Uuid: followersUUID[0]},
				{Uuid: followersUUID[1]},
			},
		})
	})
	t.Run("should_error_when_no_UUID_in_context", func(t *testing.T) {
		ctx = context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetWhoFollowPeer(ctx, &user.GetWhoFollowPeerIn{
			Uuid: uuid.New().String(),
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get user UUID from context")
	})
	t.Run("should_repo_err", func(t *testing.T) {
		userUUID := uuid.New().String()
		repoErr := errors.New("test")

		ctx = context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().GetWhoFollowPeer(ctx, userUUID).Return(nil, repoErr)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.GetWhoFollowPeer(ctx, &user.GetWhoFollowPeerIn{Uuid: userUUID})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_CheckFriendship(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	t.Run("should_ok", func(t *testing.T) {
		userUUID := uuid.New().String()
		friendUUID := uuid.New().String()
		ctx = context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, userUUID, friendUUID).Return(true, nil)

		s := &Server{dbRepo: mockDBRepo}
		res, err := s.CheckFriendship(ctx, &user.CheckFriendshipIn{Uuid: friendUUID})

		assert.NoError(t, err)
		assert.True(t, res.Succses)
	})
	t.Run("should_error_when_no_UUID_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.CheckFriendship(ctx, &user.CheckFriendshipIn{Uuid: uuid.New().String()})

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get user UUID from context")
	})
	t.Run("should_error_when_db_check_fails", func(t *testing.T) {
		userUUID := uuid.New().String()
		friendUUID := uuid.New().String()
		dbError := errors.New("database error")
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockDBRepo.EXPECT().CheckFriendship(ctx, userUUID, friendUUID).Return(false, dbError)

		s := &Server{dbRepo: mockDBRepo}
		_, err := s.CheckFriendship(ctx, &user.CheckFriendshipIn{Uuid: friendUUID})

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
	})
}

func TestServer_CreateUserPosts(t *testing.T) {
	t.Parallel()

	content := "test-content"
	userUUID := uuid.New()
	expUUID := uuid.New().String()
	msg := &user.UserPostCreated{
		UserUuid: userUUID.String(),
		PostId:   expUUID,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, config.KeyUUID, userUUID.String())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDbRepo(ctrl)
	mockKafka := NewMockUserPostCreatedProduser(ctrl)
	s := &Server{dbRepo: mockDBRepo, upcP: mockKafka}

	t.Run("create_ok", func(t *testing.T) {
		mockDBRepo.EXPECT().CreatePost(ctx, userUUID, content).Return(expUUID, nil)
		mockKafka.EXPECT().ProduceMessage(ctx, msg, userUUID).Return(nil)

		_, err := s.CreatePost(ctx, &user.CreatePostIn{Content: content})

		assert.NoError(t, err)
	})

	t.Run("create_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		_, err := s.CreatePost(ctx, &user.CreatePostIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to retrieve uuid")
	})

	t.Run("create_err", func(t *testing.T) {
		expectedErr := errors.New("get err")

		mockDBRepo.EXPECT().CreatePost(ctx, userUUID, content).Return("", expectedErr)

		_, err := s.CreatePost(ctx, &user.CreatePostIn{Content: content})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to create post: get err")
	})
}

func TestServer_GetUserPostsByIds(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)

	post1, post2 := uuid.New().String(), uuid.New().String()
	nickname1, nickname2 := "JohnDoe", "JaneDoe"
	name1, name2 := "John", "Jane"
	surname1, surname2 := "Doe", "Doe"
	avatarLink1, avatarLink2 := "avatar1", "avatar2"
	content1, content2 := "post1", "post2"
	createdAt1, createdAt2 := time.Now(), time.Now()                                                                                       //nolint:all
	editedAt1, editedAt2 := null.Time{sql.NullTime{Time: time.Now(), Valid: true}}, null.Time{sql.NullTime{Time: time.Now(), Valid: true}} //nolint:all

	t.Run("should ok", func(t *testing.T) {
		defer ctrl.Finish()

		ctx = context.WithValue(ctx, config.KeyUUID, post1)

		mockInput := &user.GetPostsByIdsIn{
			PostUuids: []string{post1, post2},
		}

		expectedPosts := model.PostInfoList{
			{
				ID:         post1,
				Nickname:   nickname1,
				Name:       name1,
				Surname:    surname1,
				AvatarLink: avatarLink1,
				Content:    content1,
				CreatedAt:  createdAt1,
				EditedAt:   editedAt1,
			},
			{
				ID:         post2,
				Nickname:   nickname2,
				Name:       name2,
				Surname:    surname2,
				AvatarLink: avatarLink2,
				Content:    content2,
				CreatedAt:  createdAt2,
				EditedAt:   editedAt2,
			},
		}
		mockDBRepo := NewMockDbRepo(ctrl)
		mockDBRepo.EXPECT().GetPostsByIds(ctx, mockInput.PostUuids).Return(&expectedPosts, nil)

		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetPostsByIds(ctx, mockInput)
		assert.NoError(t, err)
		assert.Equal(t, result, &user.GetPostsByIdsOut{Posts: (&expectedPosts).FromDTO()})
	})

	t.Run("should_return_nil_if_empty_UUID_provided", func(t *testing.T) {
		defer ctrl.Finish()

		mockInput := &user.GetPostsByIdsIn{
			PostUuids: []string{},
		}

		s := &Server{dbRepo: nil}
		_, err := s.GetPostsByIds(ctx, mockInput)
		assert.Error(t, err)
	})

	t.Run("should_return_error_if_db_fails", func(t *testing.T) {
		defer ctrl.Finish()

		ctx = context.WithValue(ctx, config.KeyUUID, post1)

		mockInput := &user.GetPostsByIdsIn{
			PostUuids: []string{post1, post2},
		}

		expectedErr := errors.New("get err")
		mockDBRepo := NewMockDbRepo(ctrl)
		mockDBRepo.EXPECT().GetPostsByIds(ctx, mockInput.PostUuids).Return(nil, expectedErr)

		s := &Server{dbRepo: mockDBRepo}
		result, err := s.GetPostsByIds(ctx, mockInput)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
