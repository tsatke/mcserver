package game

import "github.com/rs/zerolog"

type Option func(*Game)

func WithLogger(log zerolog.Logger) Option {
	return func(g *Game) {
		g.log = log
	}
}
