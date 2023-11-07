/**
 * Created by psyduck on 2022/12/14
 */

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jayaras/mdns-proxy/proxy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	_ "go.uber.org/automaxprocs"
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

	viper.SetEnvPrefix("MDNS_PROXY")
	viper.AutomaticEnv()

	rootCmd.Flags().DurationP("timeout", "t", time.Second*4, "timeout for mdns response")
	rootCmd.Flags().Uint16P("port", "p", 5345, "dns server udp port")
	rootCmd.Flags().StringP("ip", "i", "0.0.0.0", "ip address to listen on")
	rootCmd.Flags().StringP("zone", "z", "mdns.", "authoritive dns zone")
	rootCmd.Flags().BoolP("recursive", "r", true, "enable recursive resolver")
	rootCmd.Flags().StringP("upstream", "u", "192.168.1.1:53", "upstream DNS Server")

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {

		viper.BindPFlags(rootCmd.Flags())

		srv := &proxy.Server{
			Log:       zapr.NewLogger(zapLog),
			IP:        viper.GetString("ip"),
			Port:      int(viper.GetUint16("port")),
			Timeout:   viper.GetDuration("timeout"),
			Zone:      viper.GetString("zone"),
			Recusrive: viper.GetBool("recursive"),
			Upstream:  viper.GetString("upstream"),
		}

		if err = srv.ListenAndServe(); err != nil {
			return fmt.Errorf("could not start dns server: %w", err)
		}

		return nil

	}

	return rootCmd.Execute()
}
