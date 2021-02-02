package mcserver

import (
	"encoding/json"
	"sync"

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
	suite.ClosedOrEOF(netConn)
}

func (suite *ServerSuite) TestSendInvalidHandshake() {
	netConn := suite.DialServer()
	suite.DoSend(netConn, packet.IDServerboundHandshake, func(enc packet.Encoder) {
		// missing protocol version
		enc.WriteString("server address", "localhost")
		enc.WriteUshort("server port", 12345)
		enc.WriteVarInt("next state", int(packet.NextStateStatus))
	})
	// server must close connection if packet is not readable as per packet
	// definition
	suite.ClosedOrEOF(netConn)
}

func (suite *ServerSuite) TestEmptyHandshake() {
	netConn := suite.DialServer()
	suite.DoSend(netConn, packet.IDServerboundHandshake, func(enc packet.Encoder) {})
	// server must close connection if packet is not readable as per packet
	// definition
	suite.ClosedOrEOF(netConn)
}

func (suite *ServerSuite) TestSendNonHandshakePacket() {
	netConn := suite.DialServer()
	suite.DoSend(netConn, packet.IDServerboundRequest, func(enc packet.Encoder) {
		// send no fields as per packet definition
	})
	// server must close the connection because it did not receive a handshake as
	// first message
	suite.ClosedOrEOF(netConn)
}

func (suite *ServerSuite) TestMultipleConnections() {
	statusFlow := func() {
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
		suite.ClosedOrEOF(netConn)
	}
	suite.Run("mode=sequential", func() {
		cnt := 3
		for i := 0; i < cnt; i++ {
			statusFlow()
		}
	})
	suite.Run("mode=parallel", func() {
		wg := &sync.WaitGroup{}
		cnt := 3
		wg.Add(cnt)
		for i := 0; i < cnt; i++ {
			go func() {
				statusFlow()
				wg.Done()
			}()
		}
		wg.Wait()
	})
	suite.Run("mode=parallel_many", func() {
		wg := &sync.WaitGroup{}
		cnt := 100
		wg.Add(cnt)
		for i := 0; i < cnt; i++ {
			go func() {
				statusFlow()
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
