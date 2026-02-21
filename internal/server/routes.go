package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	// Public routes (no authentication required)
	r.HandlerFunc(http.MethodGet, "/", s.HelloWorldHandler)
	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)
	r.HandlerFunc(http.MethodPost, "/signup", s.SignUpHandler)
	r.HandlerFunc(http.MethodPost, "/signin", s.SignInHandler)
	r.HandlerFunc(http.MethodPost, "/refresh", s.RefreshTokenHandler)

	// Protected routes (authentication required)
	r.HandlerFunc(http.MethodGet, "/users", s.withAuth(s.getUsersHandler))
	r.HandlerFunc(http.MethodGet, "/users/:id", s.withAuth(s.getUserHandler))
	r.HandlerFunc(http.MethodPost, "/users", s.withAuth(s.createUserHandler))
	r.HandlerFunc(http.MethodPut, "/users/:id", s.withAuth(s.updateUserHandler))
	r.HandlerFunc(http.MethodDelete, "/users/:id", s.withAuth(s.deleteUserHandler))

	// Wrap all routes with CORS middleware (outermost layer)
	return s.corsMiddleware(r)
}
