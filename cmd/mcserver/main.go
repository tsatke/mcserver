package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/tsatke/mcserver"
)

var (
	// Version can be set with the Go linker.
	Version string = "main"
	// AppName is the name of this app, as displayed in the help
	// text of the root command.
	AppName = "mcserver"

	log zerolog.Logger
)

var (
	rootCmd = &cobra.Command{
		Use:  AppName,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := args[0]
			srv := mcserver.New(log, addr)
			go func() {
				<-cmd.Context().Done()
				srv.Stop()
			}()
			return srv.Start(cmd.Context())
		},
		Version: Version,
	}
)

func main() {
	log = zerolog.New(
		zerolog.ConsoleWriter{
			Out: os.Stdout,
		},
	).With().Timestamp().Logger()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// start listening for signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			log.Info().
				Msg("shutting down...")
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		log.Info().
			Msg("forced shutdown")
		os.Exit(2)
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
