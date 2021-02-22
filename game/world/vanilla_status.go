package world

type Status string

const (
	StatusEmpty               Status = "empty"
	StatusStructureStarts     Status = "structure_starts"
	StatusStructureReferences Status = "structure_references"
	StatusBiomes              Status = "biomes"
	StatusNoise               Status = "noise"
	StatusSurface             Status = "surface"
	StatusCarvers             Status = "carvers"
	StatusLiquidCarvers       Status = "liquid_carvers"
	StatusFeatures            Status = "features"
	StatusLight               Status = "light"
	StatusSpawn               Status = "spawn"
	StatusHeightmaps          Status = "heightmaps"
	StatusFull                Status = "full"
)
