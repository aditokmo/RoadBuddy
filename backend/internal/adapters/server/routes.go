package server

import (
	"backend/internal/adapters/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	corsWrapper := middleware.CORS(r)
	authenticate := middleware.JWT(s.services.Auth)

	// Auth routes
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/register", s.handlers.Auth.CreateAccount)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/login", s.handlers.Auth.Login)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/refresh", s.handlers.Auth.RefreshToken)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/logout", authenticate(s.handlers.Auth.Logout))
	// r.HandlerFunc(http.MethodGet, "/api/v1/auth/me", s.handlers.Auth.GetCurrentUser)

	// User routes
	r.HandlerFunc(http.MethodGet, "/api/v1/users", s.handlers.User.GetUsers)
	r.HandlerFunc(http.MethodGet, "/api/v1/users/:id", s.handlers.User.GetUser)

	return corsWrapper
}
