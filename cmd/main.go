package main

import (
	"net/http"
	"time"

	"github.com/k8s-container-integrity-monitor/pkg/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Info("Initialize Integrity Monitor Injector")

	mux := http.NewServeMux()
	h := handlers.New(logger)
	h.Register(mux)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	logger.Info("Start Integrity Monitor Injector webhook server")
	if err := s.ListenAndServeTLS("./ssl/k8s-webhook-injector.pem", "./ssl/k8s-webhook-injector.key"); err != nil {
		logger.WithError(err).Fatal("Failed run http server")
	}
}
