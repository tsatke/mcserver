package packet

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/tsatke/mcserver/game/chat"
)

func init() {
	RegisterPacket(StateStatus, reflect.TypeOf(ClientboundResponse{}))
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

	ResponseVersion struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	}

	ResponsePlayers struct {
		Max    int                     `json:"max"`
		Online int                     `json:"online"`
		Sample []ResponsePlayersSample `json:"sample"`
	}

	ResponsePlayersSample struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
)

func (r Response) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type ClientboundResponse struct {
	JSONResponse Response
}

func (ClientboundResponse) ID() ID       { return IDClientboundResponse }
func (ClientboundResponse) Name() string { return "Response" }

func (c ClientboundResponse) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteString("json response", c.JSONResponse.String())

	return
}
