package packet

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/tsatke/mcserver/game/chat"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ClientboundResponse{}))
}

type (
	// Response is the JSON payload of the clientbound
	// response message.
	Response struct {
		Version     ResponseVersion `json:"version"`
		Players     ResponsePlayers `json:"players"`
		Description chat.Chat       `json:"description"`
		Favicon     string          `json:"favicon,omitempty"`
	}

	// ResponseVersion holds the version and protocol of this server.
	// The version is a string, which is displayed as is in the client.
	// E.g. if a client is banned from the server, and the ban can
	// be detected in the status phase - i.e. by IP address - name
	// can be set to e.g. "forbidden", which will then be displayed
	// by the client instead of a version number.
	ResponseVersion struct {
		// Name is the name of the server version, usually similar
		// to 1.16.
		Name string `json:"name"`
		// Protocol is the protocol version of this server.
		// The client will check this value to see if it is
		// compatible with this server.
		Protocol int `json:"protocol"`
	}

	// ResponsePlayers is the player information of this server,
	// telling the client the current and max amount of players
	// on the server, as well as a sample player names that are
	// currently logged on. If the server implements some kind of
	// friendship system, the sample should prioritise befriended
	// IP addresses for the connecting IP address. As of 2020-03-01,
	// this server does not implement anything like this.
	ResponsePlayers struct {
		// Max is the maximum amount of players on this server,
		// as displayed by the client. Every value, no matter
		// positive or negative, will be displayed as it is sent.
		// I am not aware of any value that would indicate
		// to the client that the server has no upper limit of
		// connections.
		Max int `json:"max"`
		// Online is the current amount of online players on the
		// server.
		Online int `json:"online"`
		// Sample is an excerpt of online players on the server.
		// Part of this list will be displayed by the client.
		Sample []ResponsePlayersSample `json:"sample"`
	}

	// ResponsePlayersSample is a single player name and ID,
	// which is displayed by the client in the server preview
	// hover.
	ResponsePlayersSample struct {
		// Name is the name of the player.
		Name string `json:"name"`
		// ID is a unique ID of the player with the name given
		// in this struct. Does not need to relate to any game
		// mechanics, must be unique within this one packet.
		ID string `json:"id"`
	}
)

func (r Response) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// ClientboundResponse is the server's response to the empty
// ServerboundRequest packet, and is only sent in the status phase.
// The content is an aggregation of server information, that will
// be displayed by the client in the server list in the multiplayer
// menu.
type ClientboundResponse struct {
	// JSONResponse is the aggregated server data that will
	// be displayed in the multiplayer menu of the client.
	JSONResponse Response
}

// ID returns the constant packet ID.
func (ClientboundResponse) ID() ID { return IDClientboundResponse }

// Name returns the constant packet name.
func (ClientboundResponse) Name() string { return "Response" }

// EncodeInto writes this packet into the given writer.
func (c ClientboundResponse) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteString("json response", c.JSONResponse.String())

	return
}
