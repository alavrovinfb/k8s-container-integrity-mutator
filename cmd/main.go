package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ScienceSoft-Inc/k8s-container-integrity-mutator/pkg/handlers"
)

func main() {
	logger := logrus.New()
	logger.Info("Initialize Integrity Monitor Injector")

	mux := http.NewServeMux()
	h := handlers.New(logger)
	h.Register(mux)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", viper.GetInt("webhook.port")),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	logger.Info("Start Integrity Monitor Injector webhook server")
	if err := s.ListenAndServeTLS(viper.GetString("tls.cert.file"), viper.GetString("tls.key.file")); err != nil {
		logger.WithError(err).Fatal("Failed run http server")
	}
}
