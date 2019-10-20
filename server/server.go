package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Config - server config
type Config struct {
	Address string `yaml:"addr"`
	Root  string `yaml:"root"`
	Assets  string `yaml:"assets"`
}

// Server - server instance
type Server struct {
	ctx  context.Context
	lg   *logrus.Logger
	conf *Config
	mux  *chi.Mux
	serv *http.Server
}

// NewServer - creates new server instance
func NewServer(ctx context.Context, lg *logrus.Logger, conf *Config) *Server {
	mux := chi.NewMux()
	serv := &Server{
		ctx:  ctx,
		lg:   lg,
		conf: conf,
		mux:  mux,
	}
	serv.applyRouters()
	return serv
}

// Run - run current server instance
func (serv *Server) Run() error {
	serv.lg.Infof("server starting on %s", serv.conf.Address)
	srv := &http.Server{
		Addr:    serv.conf.Address,
		Handler: serv.mux,
	}
	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()
	serv.serv = srv
	select {
	case err := <-errChan:
		return err
	case <-time.After(1 * time.Second): //< if time without errors then OK
		serv.lg.Infof("started")
		return nil
	}
}

// Stop - stopping the server instance
func (serv *Server) Stop() error {
	serv.lg.Infof("server stopped")
	return serv.serv.Shutdown(serv.ctx)
}

// JSONResp - send 200 OK and response as JSON object
func (serv *Server) JSONResp(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		serv.lg.Fatalf("can't marshal object: %s", err)
	}
	if _, err := w.Write(data); err != nil {
		serv.lg.Fatalf("can't write data: %s", err)
	}
}