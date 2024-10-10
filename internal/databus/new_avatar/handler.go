package new_avatar

import (
	"context"
	"encoding/json"
	"log"

	"github.com/s21platform/metrics-lib/pkg"
	userproto "github.com/s21platform/user-proto/user-proto/new_avatar_register"
	"github.com/s21platform/user-service/internal/config"
)

// AvatarUpdateRsvMap TODO продумать куда впихнуть
type AvatarUpdateRsvMap struct {
	UUID string `json:"uuid"`
	Link string `json:"link"`
}

type Handler struct {
	dbR DBRepo
}

func New(dbR DBRepo) *Handler {
	return &Handler{dbR: dbR}
}

func convertMessage(bMessage []byte, target interface{}) error {
	err := json.Unmarshal(bMessage, target)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Handler(ctx context.Context, in []byte) {
	m := pkg.FromContext(ctx, config.KeyMetrics)

	var msg userproto.NewAvatarRegister
	err := convertMessage(in, &msg)

	if err != nil {
		m.Increment("new_avatar.error")
		log.Printf("failed to convert message: %v", err)

		return
	}

	err = h.dbR.UpdateUserAvatar(msg.Uuid, msg.Link)

	if err != nil {
		m.Increment("new_avatar.error")
		log.Printf("failed to update avatar: %v", err)

		return
	}
}
