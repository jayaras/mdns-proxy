/**
 * Created by psyduck on 2022/12/14
 */

package main

import (
	"fmt"
	"mdns-proxy/proxy"
	"os"
	"time"

	"mdns-proxy/proxy"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// just using zap but logr supports multiple frameworks
	// so it can be swapped out easily.
	zapLog, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("could not create log: %w", err)
	}

	srv := &proxy.Server{
		Log:     zapr.NewLogger(zapLog),
		IP:      "127.0.0.1",
		Port:    5678,
		Timeout: time.Second * 2,
		Zone:    "dev.",
	}

	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("could not start dns server: %w", err)
	}

	return nil
}
