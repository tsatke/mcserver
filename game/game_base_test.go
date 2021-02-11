package game

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}

type GameSuite struct {
	suite.Suite

	testdata afero.Fs
	world    afero.Fs
	game     *Game
}

func (suite *GameSuite) SetupSuite() {
	suite.testdata = afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	suite.world = afero.NewCopyOnWriteFs(afero.NewBasePathFs(suite.testdata, "maps/world01"), afero.NewMemMapFs())
	game, err := New(zerolog.Nop(), suite.world)
	suite.NoError(err)
	suite.game = game
	suite.NoError(suite.game.loadWorld())
}
