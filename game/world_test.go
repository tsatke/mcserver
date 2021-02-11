package game

import (
	"github.com/rs/zerolog"
	"github.com/spf13/afero"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

func (suite *WorldSuite) TestLoadWorld() {
	world, err := LoadWorld(
		zerolog.Nop(),
		afero.NewBasePathFs(suite.testdata, "maps/world01"),
	)
	suite.NoError(err)
	suite.Run("world.level", func() {
		world.Level.WorldGenSettings.DimensionsNBT = nil // don't compare because pointer
		suite.Equal(Level{
			BorderCenterX:        0,
			BorderCenterZ:        0,
			BorderDamagePerBlock: 0.2,
			BorderSafeZone:       5,
			BorderSize:           6e7,
			BorderSizeLerpTarget: 6e7,
			BorderSizeLerpTime:   0,
			BorderWarningBlocks:  5,
			BorderWarningTime:    15,
			DataPacks: DataPacks{
				Enabled: []string{"vanilla"},
			},
			DataVersion:      2586,
			DayTime:          3866,
			Difficulty:       1,
			DifficultyLocked: false,
			DragonFight: DragonFight{
				DragonKilled: 1,
				Gateways: []int{
					7,
					0,
					14,
					3,
					8,
					6,
					4,
					12,
					16,
					5,
					15,
					1,
					10,
					19,
					11,
					17,
					13,
					18,
					2,
					9,
				},
				PreviouslyKilled: 1,
			},
			GameRules: GameRules{
				AnnounceAdvancements:       true,
				CommandBlockOutput:         true,
				DisableElytraMovementCheck: false,
				DisableRaids:               false,
				DoDaylightCycle:            true,
				DoEntityDrops:              true,
				DoFireTick:                 true,
				DoImmediateRespawn:         false,
				DoInsomnia:                 true,
				DoLimitedCrafting:          false,
				DoMobLoot:                  true,
				DoMobSpawning:              true,
				DoPatrolSpawning:           true,
				DoTileDrops:                true,
				DoTraderSpawning:           true,
				DoWeatherCycle:             true,
				DrowningDamage:             true,
				FallDamage:                 true,
				FireDamage:                 true,
				ForgiveDeadPlayers:         true,
				KeepInventory:              false,
				LogAdminCommands:           true,
				MaxCommandChainLength:      65536,
				MaxEntityCramming:          24,
				MobGriefing:                true,
				NaturalRegeneration:        true,
				RandomTickSpeed:            3,
				ReducedDebugInfo:           false,
				SendCommandFeedback:        true,
				ShowDeathMessages:          true,
				SpawnRadius:                10,
				SpectatorsGenerateChunks:   true,
				UniversalAnger:             false,
			},
			GameType:     0,
			LastPlayed:   1610916474115,
			LevelName:    "world",
			ServerBrands: []string{"vanilla"},
			SpawnAngle:   0,
			Spawn:        voxel.V3{96, 65, -64},
			Time:         3866,
			Version: Version{
				ID:       2586,
				Name:     "1.16.5",
				Snapshot: 0,
			},
			WanderingTraderSpawnChance: 25,
			WanderingTraderSpawnDelay:  20400,
			WasModded:                  false,
			WorldGenSettings: WorldGenSettings{
				BonusChest:       false,
				DimensionsNBT:    nil, // excluded because pointer
				GenerateFeatures: true,
				Seed:             7653369538012276569,
			},
			AllowCommands:    false,
			ClearWeatherTime: 0,
			Hardcore:         false,
			Initialized:      true,
			RainTime:         174196,
			Raining:          false,
			ThunderTime:      65839,
			Thundering:       false,
		}, world.Level)
	})
	suite.Run("chunk 0,0", func() {
		reg, err := world.LoadRegion(voxel.V2{0, 0})
		suite.NoError(err)
		defer func() {
			suite.NoError(reg.Close())
		}()

		c1, err := reg.loadChunk(voxel.V2{0, 0})
		suite.Require().NoError(err)

		// unfortunately we cannot test every block in the entire chunk due to limitations
		// in the Go linker.

		suite.Run("y0_bedrock", func() {
			// check bedrock bottom layer
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					got := c1.BlockAt(voxel.V3{x, 0, z})
					suite.Equal(id.ParseID("minecraft:bedrock"), got.ID())
				}
			}
		})
	})
}
