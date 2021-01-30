package entity

//go:generate stringer -linecomment -trimprefix=Color -type=Color

type Color byte

const (
	ColorWhite Color = iota
	ColorOrange
	ColorMagenta
	ColorLightBlue
	ColorYellow
	ColorLime
	ColorPink
	ColorGray
	ColorLightGray
	ColorCyan
	ColorPurple
	ColorBlue
	ColorBrown
	ColorGreen
	ColorRed
	ColorBlack
)
