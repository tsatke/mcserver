package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tsatke/mcserver/network/packet/types"
)

func TestPacketSuite(t *testing.T) {
	suite.Run(t, new(PacketSuite))
}

type PacketSuite struct {
	suite.Suite
}

func (suite *PacketSuite) valuesReader(values ...types.Value) io.Reader {
	var buf bytes.Buffer
	for _, val := range values {
		suite.NoError(val.EncodeInto(&buf))
	}
	return &buf
}
