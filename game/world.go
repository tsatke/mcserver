package game

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/voxel"
	"github.com/tsatke/mcserver/game/worldgen"
)

type World struct {
	log zerolog.Logger
	// name is the name of the world.
	name string
	// debug indicates whether this world is in debug mode.
	// Debug mode worlds cannot be modified and have predefined blocks.
	debug bool
	// superflat indicates whether this is a superflat world.
	// Superflat worlds have different void fog and a horizon at y=0 instead of y=63.
	superflat bool

	// worldFs holds all relevant subdirectories in the savegame's directory.
	worldFs

	Level

	// generator is the world generator that is used to generate chunks.
	generator worldgen.Generator
}

type worldFs struct {
	fs         afero.Fs
	region     afero.Fs
	playerdata afero.Fs
}

func LoadWorld(log zerolog.Logger, fs afero.Fs) (*World, error) {
	world := &World{
		log: log,
		worldFs: worldFs{
			fs:         fs,
			region:     afero.NewBasePathFs(fs, "region"),
			playerdata: afero.NewBasePathFs(fs, "playerdata"),
		},
	}

	if err := world.loadLevel(); err != nil {
		return nil, fmt.Errorf("load level.dat: %w", err)
	}

	// TODO: load the respective world generator if possible
	world.generator = worldgen.NewSuperflatGenerator(
		log.With().
			Str("worldgen", "superflat").
			Logger(),
		world.WorldGenSettings.Seed,
	)

	return world, nil
}

func (w *World) LoadNBTPlayerdata(playerUUID uuid.UUID) (nbt.Tag, error) {
	playerdataFileName := fmt.Sprintf("%s.dat", playerUUID)
	playerdataFile, err := w.playerdata.Open(playerdataFileName)
	if err != nil {
		return nil, fmt.Errorf("open playerdata: %w", err)
	}

	// playerdata is gzip compressed
	decompressor, err := gzip.NewReader(playerdataFile)
	if err != nil {
		return nil, fmt.Errorf("gzip decompressor: %w", err)
	}

	dec := nbt.NewDecoder(decompressor, binary.BigEndian)
	data, err := dec.ReadTag()
	if err != nil {
		return nil, fmt.Errorf("read nbt: %w", err)
	}
	return data, nil
}

// loadRegion will return nil if there is no region created yet.
func (w *World) LoadRegion(coord voxel.V2) (*Region, error) {
	regionFileName := fmt.Sprintf("r.%d.%d.mca", coord.X, coord.Z)
	regionFile, err := w.region.Open(regionFileName)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", regionFileName, err)
	}

	if info, err := regionFile.Stat(); err != nil {
		return nil, fmt.Errorf("stat '%s': %w", regionFileName, err)
	} else if info.Size() == 0 {
		return nil, nil
	}

	return loadRegion(w.log, regionFile)
}
