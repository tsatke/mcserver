package packet

import "fmt"

type ID int

// Packet IDs.
const (
	IDServerboundHandshake                 ID = 0x00
	IDServerboundLoginStart                ID = 0x00
	IDServerboundRequest                   ID = 0x00
	IDServerboundTeleportConfirm           ID = 0x00
	IDServerboundPing                      ID = 0x01
	IDServerboundClientSettings            ID = 0x05
	IDServerboundPluginMessage             ID = 0x0B
	IDServerboundPlayerPositionAndRotation ID = 0x13

	IDClientboundResponse              ID = 0x00
	IDClientboundDisconnectLogin       ID = 0x00
	IDClientboundPong                  ID = 0x01
	IDClientboundLoginSuccess          ID = 0x02
	IDClientboundServerDifficulty      ID = 0x0D
	IDClientboundDisconnectPlay        ID = 0x19
	IDClientboundEntityStatus          ID = 0x1A
	IDClientboundChunkData             ID = 0x20
	IDClientboundUpdateLight           ID = 0x23
	IDClientboundJoinGame              ID = 0x24
	IDClientboundPlayerInfo            ID = 0x32
	IDClientboundPlayerPositionAndLook ID = 0x34
	IDClientboundHeldItemChange        ID = 0x3F
	IDClientboundUpdateViewPosition    ID = 0x40
	IDClientboundDeclareRecipes        ID = 0x5A
	IDClientboundTags                  ID = 0x5B
)

func (id ID) String() string {
	return fmt.Sprintf("0x%02x", int(id))
}
