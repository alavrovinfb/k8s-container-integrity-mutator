package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Webhook config
const (
	defaultWebhookPort = 8443
	defaultTLSCertFile = "/ssl/server.pem"
	defaultTLSKeyFile  = "/ssl/server.key"
	defaultSidecarCfg  = "/etc/sidecar/config/monitor-sidecar-config.yaml"
)

var (
	// get command line parameters
	flagPort = pflag.Int("webhook.port", defaultWebhookPort, "port of webhook server")
	flagCert = pflag.String("tls.cert.file", defaultTLSCertFile, "file containing the x509 Certificate for webhook HTTPS")
	flagKey  = pflag.String("tls.key.file", defaultTLSKeyFile, "file containing the x509 private key to --tls.key.file .")
	flagCfg  = pflag.String("sidecar.cfg.file", defaultSidecarCfg, "File containing the mutation configuration.")
)

func init() {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("can't bind flags: %v", err)
		os.Exit(2)
		return
	}
	pflag.Parse()
}
