package rest

import (
	"log"
	"net/http"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request, userIdentity string) {
	log.Println("GetApiUsersProfileUserIdentity", userIdentity)
	w.WriteHeader(http.StatusOK)
}
