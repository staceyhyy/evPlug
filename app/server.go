package app

import (
	"golive/logger"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//NewServer sets server config
func NewServer(r *mux.Router, hostAddress string) *http.Server {
	s := &http.Server{
		Addr:         hostAddress,           // configure the bind address
		Handler:      middlewareRecovery(r), // set the default handler
		ErrorLog:     logger.Error,          // set the logger for the server
		ReadTimeout:  15 * time.Second,      // max time to read request from the client
		WriteTimeout: 15 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}
	return s
}
