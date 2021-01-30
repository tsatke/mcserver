package packet

import (
	"github.com/tsatke/mcserver/network/packet/types"
)

func (suite *PacketSuite) TestServerboundHandshake_DecodeFrom() {
	rd := suite.valuesReader(
		types.NewVarInt(123),
		types.NewString("server.address.com"),
		types.NewUnsignedShort(12345),
		types.NewVarInt(2),
	)
	packet := new(ServerboundHandshake)
	suite.NoError(packet.DecodeFrom(rd))
	suite.EqualValues(123, packet.ProtocolVersion)
	suite.EqualValues("server.address.com", packet.ServerAddress)
	suite.EqualValues(12345, packet.ServerPort)
	suite.EqualValues(2, packet.NextState)
}
