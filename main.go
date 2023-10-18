/**
 * Created by psyduck on 2022/12/14
 */

package main

import (
	"fmt"
	"os"
	"time"

	"mdns-proxy/proxy"

	"github.com/spf13/cobra"

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

	rootCmd := &cobra.Command{Use: "mdns-proxy"}

	timeout := rootCmd.Flags().DurationP("timeout", "t", time.Second*2, "timeout for mdns response")
	port := rootCmd.Flags().Uint16P("port", "p", 5345, "dns server udp port")
	ip := rootCmd.Flags().StringP("ip", "i", "0.0.0.0", "ip address to listen on")
	zone := rootCmd.Flags().StringP("zone", "z", "local.", "authoritive dns zone")

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {

		srv := &proxy.Server{
			Log:     zapr.NewLogger(zapLog),
			IP:      *ip,
			Port:    int(*port),
			Timeout: *timeout,
			Zone:    *zone,
		}

		if err = srv.ListenAndServe(); err != nil {
			return fmt.Errorf("could not start dns server: %w", err)
		}

		return nil

	}

	return rootCmd.Execute()
}
