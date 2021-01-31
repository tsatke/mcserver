package packet

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPacketSuite(t *testing.T) {
	suite.Run(t, new(PacketSuite))
}

type PacketSuite struct {
	suite.Suite
}
