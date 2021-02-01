package mcserver

import (
	"encoding/json"

	"github.com/tsatke/mcserver/network/packet"
)

func (suite *ServerSuite) TestStatus() {
	netConn := suite.DialServer()
	suite.DoSend(netConn, packet.IDServerboundHandshake, func(enc packet.Encoder) {
		enc.WriteVarInt("protocol version", 754)
		enc.WriteString("server address", "localhost")
		enc.WriteUshort("server port", 12345)
		enc.WriteVarInt("next state", int(packet.NextStateStatus))
	})
	suite.DoSend(netConn, packet.IDServerboundRequest, func(enc packet.Encoder) {
		// request packet has no fields
	})
	suite.DoReceive(netConn, func(id packet.ID, dec packet.Decoder) {
		suite.Equal(packet.IDClientboundResponse, id)
		jsonResponse := dec.ReadString("json response")

		resp := packet.Response{}
		suite.NoError(json.Unmarshal([]byte(jsonResponse), &resp))
		suite.Equal("1.16.5", resp.Version.Name)
		suite.Equal(754, resp.Version.Protocol)
	})
	suite.DoSend(netConn, packet.IDServerboundPing, func(enc packet.Encoder) {
		enc.WriteLong("payload", 123456789)
	})
	suite.DoReceive(netConn, func(id packet.ID, dec packet.Decoder) {
		suite.Equal(packet.IDClientboundPong, id)
		suite.EqualValues(123456789, dec.ReadLong("payload"))
	})
	suite.EOF(netConn)
}
