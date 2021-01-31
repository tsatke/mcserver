package game

//go:generate stringer -trimprefix=Difficulty -type=Difficulty

type Difficulty int

const (
	DifficultyPeaceful Difficulty = iota
	DifficultyEasy
	DifficultyNormal
	DifficultyHard
)
