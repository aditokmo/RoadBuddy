package server

import (
	"backend/internal/adapters/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	corsWrapper := middleware.CORS(r)

	// Auth routes
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/register", s.handlers.Auth.CreateAccount)

	// User routes
	r.HandlerFunc(http.MethodGet, "/api/v1/users", s.handlers.User.GetUsers)
	r.HandlerFunc(http.MethodGet, "/api/v1/users/:id", s.handlers.User.GetUser)

	return corsWrapper
}
