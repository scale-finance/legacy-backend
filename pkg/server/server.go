package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Server struct {
	srv	*http.Server
}

// Get returns a server object
func Get() *Server {
	return &Server {
		srv: &http.Server{},
	}
}

// WithAddr appends Addr property to server
func (s *Server) WithAddr(addr string) *Server {
	s.srv.Addr = fmt.Sprintf(":%v", addr)
	return s
}

// WithHandler appends cors handler to server
func (s *Server) WithHandler(router *httprouter.Router) *Server {
	handler := cors.Default().Handler(router)
	s.srv.Handler = handler
	return s
}

// Start starts the server on specified port and using specified handlers.
// If the starting the server encounters errors it will return it.
func (s *Server) Start() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("Server missing handler")
	}

	return s.srv.ListenAndServe()
}