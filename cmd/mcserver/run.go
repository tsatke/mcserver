package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/tsatke/mcserver"
	"github.com/tsatke/mcserver/config"
)

func run(_ io.Reader, stdout io.Writer) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := zerolog.New(
		zerolog.ConsoleWriter{
			Out: stdout,
		},
	).Level(cfg.LogLevel()).
		With().
		Timestamp().
		Logger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := mcserver.New(cfg,
		mcserver.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

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
			cancel() // this will also stop the server
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		log.Info().
			Msg("forced shutdown")
		os.Exit(2)
	}()

	return srv.Start(ctx)
}

func loadConfig() (config.Config, error) {
	vp := viper.New()
	vp.SetConfigFile("config.yaml")

	if err := vp.ReadInConfig(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return config.Config{}, fmt.Errorf("read config: %w", err)
		}
		fmt.Println("can't read config file, attempting to create one")
	}

	cfg := config.New(vp)
	cfg.ApplyDefaults() // initialize default values

	if err := vp.WriteConfig(); err != nil {
		return config.Config{}, fmt.Errorf("write config: %w", err)
	}

	return cfg, nil
}
