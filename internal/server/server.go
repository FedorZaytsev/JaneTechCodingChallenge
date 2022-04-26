package server

import (
	"context"
	"jane_tech/internal/config"
	"jane_tech/internal/database"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// Server running on a provided address
type Server struct {
	serveAddress  string
	impressionUrl *url.URL

	router    *chi.Mux
	logger    *log.Logger
	connector *database.DatabaseConnector
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    s.serveAddress,
		Handler: s.router,
	}
	s.logger.Infof("Serving at %s", s.serveAddress)

	// run the server
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				s.logger.Infoln(err.Error())
			} else {
				s.logger.Errorln(err)
			}
		}
	}()
	// defer shutdown
	defer srv.Shutdown(context.Background())

	// catch SIGINT signal from user and shutdown server gracefully
	sgnl := make(chan os.Signal, 1)
	signal.Notify(sgnl,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	stop := <-sgnl
	s.logger.Infoln("Received", stop)
	s.logger.Info("Waiting for stop all jobs")
}

func NewServer(cfg *config.Config, log *log.Logger, connector *database.DatabaseConnector) (*Server, error) {
	impressionUrl, err := url.Parse(cfg.Server.ImpressionUrl)
	if err != nil {
		return nil, err
	}

	server := &Server{
		serveAddress:  cfg.Server.ServeAddress,
		impressionUrl: impressionUrl,
		logger:        log,
		connector:     connector,
	}
	server.router = NewRouter(impressionUrl, server)

	return server, nil
}
