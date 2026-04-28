package http

import (
	"backend/internal/adapters/render"
	"backend/internal/domain/user"
	"log/slog"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	service *user.Service
	logger  *slog.Logger
}

func NewUserHandler(service *user.Service, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		h.handleUserError(w, err)
		return
	}

	response := make([]userResponse, 0, len(users))
	for _, u := range users {
		response = append(response, toUserResponse(u))
	}

	render.JSON(w, http.StatusOK, response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	id = strings.TrimSpace(id)
	if id == "" {
		render.Error(w, http.StatusBadRequest, "User ID is required", "missing_id")
		return
	}

	user, err := h.service.GetUserById(r.Context(), id)
	if err != nil {
		h.handleUserError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, toUserResponse(user))
}
