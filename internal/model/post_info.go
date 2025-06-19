package model

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/guregu/null/v6"
	user_proto "github.com/s21platform/user-service/pkg/user"
)

type PostInfoList []*PostInfo

type PostInfo struct {
	ID         string    `db:"post_id"`
	Nickname   string    `db:"login"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	AvatarLink string    `db:"last_avatar_link"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	EditedAt   null.Time `db:"updated_at"`
}

func (pd *PostInfo) FromDTO() *user_proto.PostInfo {
	return &user_proto.PostInfo{
		PostUuid:   pd.ID,
		Nickname:   pd.Nickname,
		FullName:   pd.Name + " " + pd.Surname,
		AvatarLink: pd.AvatarLink,
		Content:    pd.Content,
		CreatedAt:  timestamppb.New(pd.CreatedAt),
		IsEdited:   pd.EditedAt.Valid,
	}
}

func (pdl *PostInfoList) ListFromDTO() []*user_proto.PostInfo {
	result := make([]*user_proto.PostInfo, 0, len(*pdl))

	for _, post := range *pdl {
		result = append(result, post.FromDTO())
	}

	fmt.Println(result)
	return result
}
