package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	user_proto "github.com/s21platform/user-service/pkg/user"
)

type PostInfoList []*PostInfo

type PostInfo struct {
	id         string    `db:"login"`
	Nickname   string    `db:"surname"`
	FullName   string    `db:"name"`
	AvatarLink string    `db:"last_avatar_link"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	IsEdited   time.Time `db:"edited_at"`
}

func (pd *PostInfo) FromDTO() *user_proto.PostInfo {
	result := &user_proto.PostInfo{
		PostUuid:   pd.id,
		Nickname:   pd.Nickname,
		FullName:   pd.FullName,
		AvatarLink: pd.AvatarLink,
		Content:    pd.Content,
		CreatedAt:  timestamppb.New(pd.CreatedAt),
		IsEdited:   pd.IsEdited != time.Time{},
	}

	return result
}

func (pdl *PostInfoList) ListFromDTO() []*user_proto.PostInfo {
	result := make([]*user_proto.PostInfo, 0, len(*pdl))

	for _, post := range *pdl {
		result = append(result, post.FromDTO())
	}

	return result
}
