package world

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"

	"github.com/tsatke/mcserver/game/voxel"
)

type vanillaWorld struct {
	fs afero.Fs

	regionsLock sync.Mutex
	regions     map[voxel.V2]*vanillaRegion

	loadedChunksLock sync.Mutex
	loadedChunks     map[voxel.V2]*vanillaChunk
}

func LoadVanilla(fs afero.Fs) (World, error) {
	w := newVanillaWorld(fs)

	if err := w.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return w, nil
}

// newVanillaWorld reads a vanilla minecraft worlds from the given file system.
func newVanillaWorld(fs afero.Fs) *vanillaWorld {
	return &vanillaWorld{
		fs:           fs,
		loadedChunks: map[voxel.V2]*vanillaChunk{},
	}
}

func (w *vanillaWorld) Chunk(v2 voxel.V2) (Chunk, error) {
	w.loadedChunksLock.Lock()
	defer w.loadedChunksLock.Unlock()

	if ch, ok := w.loadedChunks[v2]; ok {
		return ch, nil
	}
	loaded, err := w.readChunk(v2)
	if errors.Is(err, ErrChunkNotExist) { // check if chunk needs to be generated
		return w.generateChunk(v2), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read chunk: %w", err)
	}

	w.loadedChunks[v2] = loaded

	return loaded, nil
}

func (w *vanillaWorld) IsChunkLoaded(v2 voxel.V2) bool {
	w.loadedChunksLock.Lock()
	defer w.loadedChunksLock.Unlock()
	_, ok := w.loadedChunks[v2]
	return ok
}

func (w *vanillaWorld) Unload(v2 voxel.V2) {
	delete(w.loadedChunks, v2)
}

func (w *vanillaWorld) Seed() int64 {
	panic("implement me")
}

func (w *vanillaWorld) validate() error {
	return nil
}

func (w *vanillaWorld) region(v2 voxel.V2) (*vanillaRegion, error) {
	w.regionsLock.Lock()
	defer w.regionsLock.Unlock()

	if reg, ok := w.regions[v2]; ok {
		return reg, nil
	}

	regionFileName := fmt.Sprintf("r.%d.%d.mca", v2.X, v2.Z)
	regionFile, err := w.fs.Open(filepath.Join("region", regionFileName))
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", regionFileName, err)
	}

	region, err := loadVanillaRegion(regionFile)
	if err != nil {
		return nil, err
	}
	w.regions[v2] = region

	return region, nil
}

// readChunk reads the chunk with the given chunk from the disk. If the chunk
// does not exist, an error will be returned. The returned chunk will NOT be
// considered loaded.
func (w *vanillaWorld) readChunk(v2 voxel.V2) (*vanillaChunk, error) {
	regionCoord := voxel.V2{v2.X >> 5, v2.Z >> 5}
	region, err := w.region(regionCoord)
	if err != nil {
		return nil, fmt.Errorf("region: %w", err)
	}

	loaded, err := region.loadChunk(v2)
	if err != nil {
		return nil, fmt.Errorf("load chunk: %w", err)
	}
	return loaded, nil
}

func (w *vanillaWorld) generateChunk(v2 voxel.V2) Chunk {
	panic("implement me")
}
