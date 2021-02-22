package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldSuite))
}

type WorldSuite struct {
	suite.Suite
}
