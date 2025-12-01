package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"

	logger_lib "github.com/s21platform/logger-lib"
	optionhub_lib "github.com/s21platform/optionhub-lib"
	"github.com/s21platform/user-service/pkg/user"

	api "github.com/s21platform/user-service/internal/generated"
	"github.com/s21platform/user-service/internal/model"
)

type Handler struct {
	dbRepo          DbRepo
	optionHubClient OptionhubClient
	postCreatedPrd  PostCreatedProducer
}

func New(dbR DbRepo, ohC OptionhubClient, postCreatedPrd PostCreatedProducer) *Handler {
	return &Handler{
		dbRepo:          dbR,
		optionHubClient: ohC,
		postCreatedPrd:  postCreatedPrd,
	}
}

func (h *Handler) checkPostOwnership(ctx context.Context, postID string, ownerUUID uuid.UUID) error {
	return h.dbRepo.CheckPostOwnership(ctx, postID, ownerUUID)
}

func (h *Handler) MyPersonality(w http.ResponseWriter, r *http.Request, params api.MyPersonalityParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := logger_lib.WithField(r.Context(), "user_uuid", params.XUserUuid)

	personality, err := h.dbRepo.GetPersonalityByUuid(ctx, params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get personality data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	options, err := h.optionHubClient.GetAttributesMeta(ctx, model.PersonalityForm)
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
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)

	if len(params.AttributeIds) == 0 {
		logger_lib.Error(ctx, "attribute_ids parameter is empty")
		resolveError(&w, http.StatusBadRequest)
		return
	}

	userAttributes, err := h.dbRepo.GetUserAttributesByUuid(ctx, params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get user attributes data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	attributeIds := make([]model.Attribute, len(params.AttributeIds))
	for i, id := range params.AttributeIds {
		attributeIds[i] = model.Attribute(id)
	}

	options, err := h.optionHubClient.GetAttributesMeta(ctx, attributeIds)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to get attributes metadata")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	result := mapUserAttributesToAttributeItems(userAttributes, options, params.AttributeIds)

	response := api.UserAttributesResponse{
		Data: result,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request, params api.UpdateProfileParams) {
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to read body")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	var body api.AttributesValues
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to unmarshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	res, err := optionhub_lib.ParseAttributes(ctx, body.Attributes)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to parse data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	data := mapAttributeToFields(ctx, res)
	err = h.dbRepo.UpdateProfile(ctx, data, params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to update profile")
		resolveError(&w, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request, params api.CreatePostParams) {
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)

	var body api.CreatePost
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to unmarshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	ownerUUID, err := uuid.Parse(params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to parse user uuid")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	newPostUUID, err := h.dbRepo.CreatePost(ctx, ownerUUID, body.Content)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to create post")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	msg := &user.UserPostCreated{
		UserUuid: params.XUserUuid,
		PostId:   newPostUUID,
	}

	err = h.postCreatedPrd.ProduceMessage(ctx, msg, ownerUUID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to create post")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	response := api.CreatePostResponse{
		PostId: newPostUUID,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request, postId string, params api.EditPostParams) {
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)
	ctx = logger_lib.WithField(ctx, "post_id", postId)

	var body api.EditPost
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to unmarshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	ownerUUID, err := uuid.Parse(params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to parse user uuid")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	err = h.checkPostOwnership(ctx, postId, ownerUUID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to check post ownership")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	updatedContent, err := h.dbRepo.EditPost(ctx, postId, ownerUUID, body.Content)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to edit post")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	response := api.EditPostResponse{
		PostId:  postId,
		Content: updatedContent,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request, postId string, params api.DeletePostParams) {
	ctx := logger_lib.WithUserUuid(r.Context(), params.XUserUuid)
	ctx = logger_lib.WithField(ctx, "post_id", postId)

	ownerUUID, err := uuid.Parse(params.XUserUuid)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to parse user uuid")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	err = h.checkPostOwnership(ctx, postId, ownerUUID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to check post ownership")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	err = h.dbRepo.DeletePost(ctx, postId, ownerUUID)
	if err != nil {
		logger_lib.Error(logger_lib.WithError(ctx, err), "failed to delete post")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
