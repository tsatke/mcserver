package game

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldSuite))
}

type WorldSuite struct {
	suite.Suite

	testdata afero.Fs
}

func (suite *WorldSuite) SetupSuite() {
	suite.testdata = afero.NewCopyOnWriteFs(afero.NewBasePathFs(afero.NewOsFs(), "testdata"), afero.NewMemMapFs())
}
