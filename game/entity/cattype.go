package entity

//go:generate stringer -linecomment -type=CatType

type CatType int

const (
	CatTypeTabby            CatType = iota // Tabby
	CatTypeTuxedo                          // Tuxedo
	CatTypeRed                             // Red
	CatTypeSiamese                         // Siamese
	CatTypeBritishShorthair                // British Shorthair
	CatTypeCalico                          // Calico
	CatTypePersian                         // Persian
	CatTypeRagdoll                         // Ragdoll
	CatTypeWhite                           // White
	CatTypeJellie                          // Jellie
	CatTypeBlack                           // Black
)
