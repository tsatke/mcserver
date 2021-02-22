package game

import (
	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/entity"
)

func (suite *GameSuite) TestLoadPlayerEntity() {
	suite.T().Skip("LoadPlayerEntity unimplemented")

	playerUUID, err := uuid.Parse("eac89a9e-9471-4ce4-ab2e-1b2a725a57fd")
	suite.NoError(err)
	player := &Player{
		Player: &entity.Player{
			Mob: entity.Mob{
				Data: entity.Data{UUID: playerUUID}, // uuid should be enough to load player data
			},
		},
	}
	suite.NoError(suite.game.loadPlayerEntity(player))
	suite.EqualValues(playerUUID, player.UUID) // UUID must not have been modified

	e := player.Player
	suite.EqualValues(0, e.AbsorptionAmount)
	suite.EqualValues(300, e.Air)
	// TODO: check Attributes
	// TODO: check Brain
	suite.EqualValues(2586, e.DataVersion)
	suite.EqualValues(0, e.DeathTime)
	suite.EqualValues("minecraft:overworld", e.Dimension)
	suite.Len(e.EnderItems, 0)
	suite.EqualValues(0, e.FallDistance)
	suite.False(e.FallFlying)
	suite.EqualValues(-20, e.Fire)
	suite.EqualValues(20, e.Health)
	suite.EqualValues(0, e.HurtByTimestamp)
	suite.EqualValues(0, e.HurtTime)
	// TODO: check Inventory
	suite.False(e.Invulnerable)
	suite.EqualValues([3]float64{0, -0.0784000015258789, 0}, e.Motion) // that's just what's stored in the file
	suite.True(e.OnGround)
	suite.EqualValues(0, e.PortalCooldown)
	suite.EqualValues([3]float64{0.5, 68, 0.5}, e.Pos)
	suite.EqualValues([2]float32{7.499985694885254, 17.550012588500977}, e.Rotation)
	suite.EqualValues(0, e.Score)
	suite.EqualValues(0, e.SelectedItemSlot)
	suite.EqualValues(0, e.SleepTimer)
	suite.EqualValues(0, e.XPLevel)
	suite.EqualValues(0, e.XPPercentage)
	suite.EqualValues(0, e.XPSeed)
	suite.EqualValues(0, e.XPTotal)
	suite.EqualValues(0.05000000074505806, e.Abilities.FlySpeed)
	suite.False(e.Abilities.Flying)
	suite.True(e.Abilities.InstaBuild)
	suite.True(e.Abilities.Invulnerable)
	suite.True(e.Abilities.MayBuild)
	suite.True(e.Abilities.MayFly)
	suite.EqualValues(0.10000000149011612, e.Abilities.WalkSpeed)
}
