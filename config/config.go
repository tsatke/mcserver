package config

import (
	"net"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Config keys according to
// https://crushit.atlassian.net/wiki/x/AQBiCQ
const (
	KeyServerAddress = "server.address"
	KeyServerPort    = "server.port"
	KeyGameWorld     = "game.world"
	KeyLogLevel      = "log.level"
)

type Config struct {
	vp *viper.Viper
}

func New(vp *viper.Viper) Config {
	return Config{
		vp: vp,
	}
}

func (c Config) ApplyDefaults() {
	/*
		Initialize default values.
		If you change anything, this must be documented on that page.
	*/
	c.vp.SetDefault(KeyServerAddress, "localhost")
	c.vp.SetDefault(KeyServerPort, 25565)

	c.vp.SetDefault(KeyGameWorld, "world")

	c.vp.SetDefault(KeyLogLevel, "info")
}

func (c Config) LogLevel() zerolog.Level {
	lvl := zerolog.InfoLevel
	switch c.vp.GetString(KeyLogLevel) {
	case "off":
		lvl = zerolog.NoLevel
	case "error":
		lvl = zerolog.ErrorLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "debug":
		lvl = zerolog.DebugLevel
	case "trace":
		lvl = zerolog.TraceLevel
	}
	return lvl
}

func (c Config) ServerAddr() string {
	return net.JoinHostPort(c.vp.GetString(KeyServerAddress), c.vp.GetString(KeyServerPort))
}

func (c Config) GameWorld() string {
	return c.vp.GetString(KeyGameWorld)
}
