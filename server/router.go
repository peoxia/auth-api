package server

import "github.com/peoxia/auth-api/server/handler"

const v1API string = "/api/v1"

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/_healthz", handler.Healthz).Methods("GET").Name("Health")
	s.Router.HandleFunc("/", handler.Front).Methods("GET").Name("Front")

	// API routes
	api := s.Router.PathPrefix(v1API).Subrouter()
	api.HandleFunc("/login", handler.Login(s.GoogleClient)).Methods("GET").Name("Login")
	api.HandleFunc("/login/callback", handler.LoginCallback(*s.Config, s.GoogleClient, s.MongoDB)).Methods("GET").Name("LoginCallback")

	api.HandleFunc("/users/me", handler.CurrentUser(*s.Config, s.MongoDB)).Methods("GET").Name("CurrentUser")
	api.HandleFunc("/users/me", handler.UpdateUser(*s.Config, s.MongoDB)).Methods("POST").Name("UpdateUser")
	api.HandleFunc("/users/me", handler.DeleteUser(*s.Config, s.MongoDB)).Methods("DELETE").Name("DeleteUser")
}
