package server

import (
	"backend/internal/adapters/handler"
	"backend/internal/adapters/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	corsWrapper := middleware.CORS(r)

	// Auth routes
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/create-account", handler.CreateAccount)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/login", handler.Login)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/logout", handler.Logout)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/forgot-password", handler.ForgotPassword)
	r.HandlerFunc(http.MethodPost, "/api/v1/auth/verify-email", handler.VerifyEmail)

	return corsWrapper
}