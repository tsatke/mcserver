package game

//go:generate stringer -trimprefix=Gamemode -type=Gamemode

type Gamemode int

const (
	GamemodeSurvival Gamemode = iota
	GamemodeCreative
	GamemodeAdventure
	GamemodeSpectator
)
