package packet

import "bytes"

func (suite *PacketSuite) TestServerboundHandshake_DecodeFrom() {
	var buf bytes.Buffer
	enc := Encoder{&buf}
	enc.WriteVarInt("protocol version", 123)
	enc.WriteString("server address", "my.server.com")
	enc.WriteUshort("server port", 17)
	enc.WriteVarInt("next state", 1)

	var p ServerboundHandshake
	suite.NoError(p.DecodeFrom(&buf))
	suite.EqualValues(123, p.ProtocolVersion)
	suite.EqualValues("my.server.com", p.ServerAddress)
	suite.EqualValues(17, p.ServerPort)
	suite.EqualValues(NextStateStatus, p.NextState)
}
