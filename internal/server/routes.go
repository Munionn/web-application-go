package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	// Wrap all routes with CORS middleware
	corsWrapper := s.corsMiddleware(r)

	r.HandlerFunc(http.MethodGet, "/", s.HelloWorldHandler)
	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)

	// Example GORM routes (uncomment to use)
	r.HandlerFunc(http.MethodGet, "/users", s.getUsersHandler)
	r.HandlerFunc(http.MethodGet, "/users/:id", s.getUserHandler)
	r.HandlerFunc(http.MethodPost, "/users", s.createUserHandler)
	r.HandlerFunc(http.MethodPut, "/users/:id", s.updateUserHandler)
	r.HandlerFunc(http.MethodDelete, "/users/:id", s.deleteUserHandler)

	r.HandlerFunc(http.MethodPost, "/signup", s.SignUpHandler)
	r.HandlerFunc(http.MethodPost, "/signin", s.SignInHandler)
	return corsWrapper
}
