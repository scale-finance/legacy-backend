package server

import (
	"errors"
	"net/http"
	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Server struct {
	srv	*http.Server
}

// returns a server object
func Get() *Server {
	return &Server {
		srv: &http.Server{},
	}
}

// appends Addr property to server
func (s *Server) WithAddr(addr string) *Server {
	s.srv.Addr = fmt.Sprintf(":%v", addr)
	return s
}

// appends cors handler to server
func (s *Server) WithHandler(router *httprouter.Router) *Server {
	handler := cors.Default().Handler(router)
	s.srv.Handler = handler
	return s
}

// Starts the server on specified port and using specified handlers.
// If the starting the server encouters errors it will return it.
func (s *Server) Start() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("Server missing handler")
	}

	return s.srv.ListenAndServe()
}