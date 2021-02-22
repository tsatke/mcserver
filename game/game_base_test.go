package game

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/tsatke/mcserver/game/world"
)

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}

type GameSuite struct {
	suite.Suite

	testdata afero.Fs
	world    world.World
	game     *Game
}

func (suite *GameSuite) SetupSuite() {
	suite.testdata = afero.NewBasePathFs(afero.NewOsFs(), "testdata")

	worldFs := afero.NewCopyOnWriteFs(afero.NewBasePathFs(suite.testdata, "maps/world01"), afero.NewMemMapFs())
	world, err := world.LoadVanilla(worldFs)
	suite.NoError(err)
	suite.world = world

	game, err := New(suite.world)
	suite.NoError(err)
	suite.game = game
}
