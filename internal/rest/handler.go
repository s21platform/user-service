package rest

import (
	"encoding/json"
	optionhub_lib "github.com/s21platform/optionhub-lib"
	"io"
	"log"
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

func (h *Handler) GetUserAttributes(w http.ResponseWriter, r *http.Request, params api.GetUserAttributesParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := logger_lib.WithField(r.Context(), "user_uuid", params.XUserUuid)

	// Валидируем параметры
	if len(params.AttributeIds) == 0 {
		logger_lib.Error(ctx, "attribute_ids parameter is empty")
		resolveError(&w, http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя из БД
	userAttributes, err := h.dbR.GetUserAttributesByUuid(ctx, params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get user attributes data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	// Преобразуем IDs в модельные атрибуты
	attributeIds := make([]model.Attribute, len(params.AttributeIds))
	for i, id := range params.AttributeIds {
		attributeIds[i] = model.Attribute(id)
	}

	// Получаем метаданные атрибутов
	options, err := h.ohC.GetAttributesMeta(ctx, attributeIds)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get attributes metadata")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	// Маппируем данные в ответ
	result := mapUserAttributesToAttributeItems(userAttributes, options, params.AttributeIds)

	response := api.UserAttributesResponse{
		Data: result,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(resp)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request, params api.UpdateProfileParams) {
	ctx := r.Context()
	t, err := io.ReadAll(r.Body)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to read body")
		resolveError(&w, http.StatusInternalServerError)
		return
	}
	var body api.AttributesValues
	err = json.Unmarshal(t, &body)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to unmarshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}
	res, err := optionhub_lib.ParseAttributes(ctx, body.Attributes)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to unmarshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	data := mapAttributeToFields(ctx, res)
	err = h.dbR.UpdateProfile(ctx, data, params.XUserUuid)
	if err != nil {
		log.Println(err)
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to update profile")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	return
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
