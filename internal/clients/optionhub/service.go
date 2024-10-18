package optoinhub

import (
	"context"
	"fmt"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

type Handle struct {
	client optionhubproto.OptionhubServiceClient
}

func (h *Handle) GetOs(ctx context.Context, id *int64) (string, error) {
	if id == nil {
		return "", fmt.Errorf("no os id for this user")
	}

	os, err := h.client.GetOsById(ctx, &optionhubproto.GetByIdIn{Id: *id})
	if err != nil {
		return "", err
	}
	return os.Value, nil
}
