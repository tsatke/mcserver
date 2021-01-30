package packet

import "fmt"

type ID int

// Packet IDs.
const (
	IDServerboundHandshake      ID = 0x00
	IDServerboundLoginStart     ID = 0x00
	IDServerboundRequest        ID = 0x00
	IDServerboundPing           ID = 0x01
	IDServerboundPluginMessage  ID = 0x0B
	IDServerboundClientSettings ID = 0x05

	IDClientboundResponse         ID = 0x00
	IDClientboundPong             ID = 0x01
	IDClientboundDisconnectLogin  ID = 0x00
	IDClientboundDisconnectPlay   ID = 0x19
	IDClientboundLoginSuccess     ID = 0x02
	IDClientboundJoinGame         ID = 0x24
	IDClientboundServerDifficulty ID = 0x0D
	IDClientboundHeldItemChange   ID = 0x3F
	IDClientboundDeclareRecipes   ID = 0x5A
)

func (id ID) String() string {
	return fmt.Sprintf("0x%02x", int(id))
}
