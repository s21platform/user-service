package rest

import (
	"encoding/json"
	"net/http"

	logger_lib "github.com/s21platform/logger-lib"

	api "github.com/s21platform/user-service/internal/generated"
	"github.com/s21platform/user-service/internal/model"
)

type Handler struct {
	dbR DbRepo
	ohC OptionhubClient
}

func New(dbR DbRepo, ohC OptionhubClient) *Handler {
	return &Handler{
		dbR: dbR,
		ohC: ohC,
	}
}

func (h *Handler) MyPersonality(w http.ResponseWriter, r *http.Request, params api.MyPersonalityParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := logger_lib.WithField(r.Context(), "user_uuid", params.XUserUuid)

	personality, err := h.dbR.GetPersonalityByUuid(ctx, params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get personality data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	options, err := h.ohC.GetAttributesMeta(ctx, model.PersonalityForm)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get options metadata")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	result := mapPersonalityToProfileItems(personality, options)

	resp, err := json.Marshal(result)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(resp)
}

func resolveError(w *http.ResponseWriter, status int) {
	var message string
	switch status {
	case http.StatusBadRequest:
		message = "Произошла ошибка, попробуйте перезагрузить страницу"
	case http.StatusUnauthorized:
		message = "Вы не авторизованы для этого действия"
	case http.StatusNotFound:
		message = "Страница не найдена"
	case http.StatusInternalServerError:
		message = "У нас что-то сломалось, но мы уже чиним!"
	default:
		message = "У нас что-то сломалось, но мы уже чиним!"
	}

	body, _ := json.Marshal(api.Forbidden{
		Message: message,
	})
	(*w).WriteHeader(status)
	_, _ = (*w).Write(body)
}
