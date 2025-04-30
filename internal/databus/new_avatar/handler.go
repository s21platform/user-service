package new_avatar

import (
	"context"
	"encoding/json"
	"log"

	"github.com/s21platform/avatar-service/pkg/avatar"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/user-service/internal/config"
)

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

func (h *Handler) Handler(ctx context.Context, in []byte) error {
	m := pkg.FromContext(ctx, config.KeyMetrics)

	var msg avatar.NewAvatarRegister
	err := convertMessage(in, &msg)

	if err != nil {
		m.Increment("new_avatar.error")
		log.Printf("failed to convert message: %v", err)
		return err
	}

	err = h.dbR.UpdateUserAvatar(msg.Uuid, msg.Link)

	if err != nil {
		m.Increment("new_avatar.error")
		log.Printf("failed to update avatar: %v", err)
		return err
	}
	return nil
}
