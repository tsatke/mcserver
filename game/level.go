package game

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/voxel"
)

type (
	Level struct {
		BorderCenterX              float64
		BorderCenterZ              float64
		BorderDamagePerBlock       float64
		BorderSafeZone             float64
		BorderSize                 float64
		BorderSizeLerpTarget       float64
		BorderSizeLerpTime         int64
		BorderWarningBlocks        float64
		BorderWarningTime          float64
		CustomBossEvents           interface{} // not known
		DataPacks                  DataPacks
		DataVersion                int
		DayTime                    int64
		Difficulty                 int8
		DifficultyLocked           bool
		DragonFight                DragonFight
		GameRules                  GameRules
		GameType                   int
		LastPlayed                 int64
		LevelName                  string
		ScheduledEvents            interface{} // not known
		ServerBrands               []string
		SpawnAngle                 float32
		Spawn                      voxel.V3
		Time                       int64
		Version                    Version
		WanderingTraderSpawnChance int
		WanderingTraderSpawnDelay  int
		WasModded                  bool
		WorldGenSettings           WorldGenSettings
		AllowCommands              bool
		ClearWeatherTime           int
		Hardcore                   bool
		Initialized                bool
		RainTime                   int
		Raining                    bool
		ThunderTime                int
		Thundering                 bool
	}

	WorldGenSettings struct {
		BonusChest bool
		// DimensionsNBT is the NBT tag that is sent with the JoinGame message to the client.
		DimensionsNBT nbt.Tag
		// GenerateFeatures indicates whether structures are generated.
		GenerateFeatures bool
		// Seed is the world's seed.
		Seed int64
	}

	Version struct {
		ID       int
		Name     string
		Snapshot int8
	}

	DataPacks struct {
		Disabled []string
		Enabled  []string
	}

	DragonFight struct {
		DragonKilled     int8 // may be a bool
		Gateways         []int
		PreviouslyKilled int8 // may be a bool
	}

	GameRules struct {
		AnnounceAdvancements       bool
		CommandBlockOutput         bool
		DisableElytraMovementCheck bool
		DisableRaids               bool
		DoDaylightCycle            bool
		DoEntityDrops              bool
		DoFireTick                 bool
		DoImmediateRespawn         bool
		DoInsomnia                 bool
		DoLimitedCrafting          bool
		DoMobLoot                  bool
		DoMobSpawning              bool
		DoPatrolSpawning           bool
		DoTileDrops                bool
		DoTraderSpawning           bool
		DoWeatherCycle             bool
		DrowningDamage             bool
		FallDamage                 bool
		FireDamage                 bool
		ForgiveDeadPlayers         bool
		KeepInventory              bool
		LogAdminCommands           bool
		MaxCommandChainLength      int
		MaxEntityCramming          int
		MobGriefing                bool
		NaturalRegeneration        bool
		RandomTickSpeed            int
		ReducedDebugInfo           bool
		SendCommandFeedback        bool
		ShowDeathMessages          bool
		SpawnRadius                int
		SpectatorsGenerateChunks   bool
		UniversalAnger             bool
	}
)

func (w *World) loadLevel() (err error) {
	levelDat, err := w.worldFs.fs.Open("level.dat")
	if err != nil {
		return fmt.Errorf("open level.dat: %w", err)
	}
	defer func() {
		_ = levelDat.Close()
	}()

	uncompressor, err := gzip.NewReader(levelDat)
	if err != nil {
		return fmt.Errorf("create gzip decompressor: %w", err)
	}

	dec := nbt.NewDecoder(uncompressor, binary.BigEndian)
	tag, err := dec.ReadTag()
	if err != nil {
		return fmt.Errorf("read root tag: %w", err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			if recErr, ok := rec.(error); ok {
				err = recErr
			} else {
				panic(rec)
			}
		}
	}()
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	mapper := nbt.NewSimpleMapper(tag)
	level := &w.Level

	must(mapper.MapDouble("Data.BorderCenterX", &level.BorderCenterX))
	must(mapper.MapDouble("Data.BorderCenterZ", &level.BorderCenterZ))
	must(mapper.MapDouble("Data.BorderDamagePerBlock", &level.BorderDamagePerBlock))
	must(mapper.MapDouble("Data.BorderSafeZone", &level.BorderSafeZone))
	must(mapper.MapDouble("Data.BorderSize", &level.BorderSize))
	must(mapper.MapDouble("Data.BorderSizeLerpTarget", &level.BorderSizeLerpTarget))
	must(mapper.MapLong("Data.BorderSizeLerpTime", &level.BorderSizeLerpTime))
	must(mapper.MapDouble("Data.BorderWarningBlocks", &level.BorderWarningBlocks))
	must(mapper.MapDouble("Data.BorderWarningTime", &level.BorderWarningTime))
	// TODO: custom boss events
	must(mapper.MapCustom("Data.DataPacks", func(dataPacks nbt.Tag) error {
		var disabled, enabled []string
		mapper := nbt.NewSimpleMapper(dataPacks)
		must(mapper.MapList(
			"Disabled",
			func(size int) {
				if size > 0 {
					disabled = make([]string, size)
				}
			},
			func(i int, mapper nbt.Mapper) error {
				return mapper.MapString("", &disabled[i])
			},
		))
		must(mapper.MapList(
			"Enabled",
			func(size int) {
				if size > 0 {
					enabled = make([]string, size)
				}
			},
			func(i int, mapper nbt.Mapper) error {
				return mapper.MapString("", &enabled[i])
			},
		))
		level.DataPacks.Disabled = disabled
		level.DataPacks.Enabled = enabled
		return nil
	}))
	must(mapper.MapInt("Data.DataVersion", &level.DataVersion))
	must(mapper.MapLong("Data.DayTime", &level.DayTime))
	must(mapper.MapByte("Data.Difficulty", &level.Difficulty))
	must(mapper.MapCustom("Data.DifficultyLocked", byteToBool(&level.DifficultyLocked)))
	must(mapper.MapCustom("Data.DragonFight", func(dragonFight nbt.Tag) error {
		var gateways []int
		mapper := nbt.NewSimpleMapper(dragonFight)
		must(mapper.MapList(
			"Gateways",
			func(size int) { gateways = make([]int, size) },
			func(i int, mapper nbt.Mapper) error {
				return mapper.MapInt("", &gateways[i])
			},
		))
		must(mapper.MapByte("DragonKilled", &level.DragonFight.DragonKilled))
		must(mapper.MapByte("PreviouslyKilled", &level.DragonFight.PreviouslyKilled))
		level.DragonFight.Gateways = gateways
		return nil
	}))
	must(mapper.MapCustom("Data.GameRules.announceAdvancements", stringToBool(&level.GameRules.AnnounceAdvancements)))
	must(mapper.MapCustom("Data.GameRules.commandBlockOutput", stringToBool(&level.GameRules.CommandBlockOutput)))
	must(mapper.MapCustom("Data.GameRules.disableElytraMovementCheck", stringToBool(&level.GameRules.DisableElytraMovementCheck)))
	must(mapper.MapCustom("Data.GameRules.disableRaids", stringToBool(&level.GameRules.DisableRaids)))
	must(mapper.MapCustom("Data.GameRules.doDaylightCycle", stringToBool(&level.GameRules.DoDaylightCycle)))
	must(mapper.MapCustom("Data.GameRules.doEntityDrops", stringToBool(&level.GameRules.DoEntityDrops)))
	must(mapper.MapCustom("Data.GameRules.doFireTick", stringToBool(&level.GameRules.DoFireTick)))
	must(mapper.MapCustom("Data.GameRules.doImmediateRespawn", stringToBool(&level.GameRules.DoImmediateRespawn)))
	must(mapper.MapCustom("Data.GameRules.doInsomnia", stringToBool(&level.GameRules.DoInsomnia)))
	must(mapper.MapCustom("Data.GameRules.doLimitedCrafting", stringToBool(&level.GameRules.DoLimitedCrafting)))
	must(mapper.MapCustom("Data.GameRules.doMobLoot", stringToBool(&level.GameRules.DoMobLoot)))
	must(mapper.MapCustom("Data.GameRules.doMobSpawning", stringToBool(&level.GameRules.DoMobSpawning)))
	must(mapper.MapCustom("Data.GameRules.doPatrolSpawning", stringToBool(&level.GameRules.DoPatrolSpawning)))
	must(mapper.MapCustom("Data.GameRules.doTileDrops", stringToBool(&level.GameRules.DoTileDrops)))
	must(mapper.MapCustom("Data.GameRules.doTraderSpawning", stringToBool(&level.GameRules.DoTraderSpawning)))
	must(mapper.MapCustom("Data.GameRules.doWeatherCycle", stringToBool(&level.GameRules.DoWeatherCycle)))
	must(mapper.MapCustom("Data.GameRules.drowningDamage", stringToBool(&level.GameRules.DrowningDamage)))
	must(mapper.MapCustom("Data.GameRules.fallDamage", stringToBool(&level.GameRules.FallDamage)))
	must(mapper.MapCustom("Data.GameRules.fireDamage", stringToBool(&level.GameRules.FireDamage)))
	must(mapper.MapCustom("Data.GameRules.forgiveDeadPlayers", stringToBool(&level.GameRules.ForgiveDeadPlayers)))
	must(mapper.MapCustom("Data.GameRules.keepInventory", stringToBool(&level.GameRules.KeepInventory)))
	must(mapper.MapCustom("Data.GameRules.logAdminCommands", stringToBool(&level.GameRules.LogAdminCommands)))
	must(mapper.MapCustom("Data.GameRules.maxCommandChainLength", stringToInt(&level.GameRules.MaxCommandChainLength)))
	must(mapper.MapCustom("Data.GameRules.maxEntityCramming", stringToInt(&level.GameRules.MaxEntityCramming)))
	must(mapper.MapCustom("Data.GameRules.mobGriefing", stringToBool(&level.GameRules.MobGriefing)))
	must(mapper.MapCustom("Data.GameRules.naturalRegeneration", stringToBool(&level.GameRules.NaturalRegeneration)))
	must(mapper.MapCustom("Data.GameRules.randomTickSpeed", stringToInt(&level.GameRules.RandomTickSpeed)))
	must(mapper.MapCustom("Data.GameRules.reducedDebugInfo", stringToBool(&level.GameRules.ReducedDebugInfo)))
	must(mapper.MapCustom("Data.GameRules.sendCommandFeedback", stringToBool(&level.GameRules.SendCommandFeedback)))
	must(mapper.MapCustom("Data.GameRules.showDeathMessages", stringToBool(&level.GameRules.ShowDeathMessages)))
	must(mapper.MapCustom("Data.GameRules.spawnRadius", stringToInt(&level.GameRules.SpawnRadius)))
	must(mapper.MapCustom("Data.GameRules.spectatorsGenerateChunks", stringToBool(&level.GameRules.SpectatorsGenerateChunks)))
	must(mapper.MapCustom("Data.GameRules.universalAnger", stringToBool(&level.GameRules.UniversalAnger)))
	must(mapper.MapInt("Data.GameType", &level.GameType))
	must(mapper.MapLong("Data.LastPlayed", &level.LastPlayed))
	must(mapper.MapString("Data.LevelName", &level.LevelName))
	// TODO: scheduled events
	must(mapper.MapList(
		"Data.ServerBrands",
		func(size int) { level.ServerBrands = make([]string, size) },
		func(i int, mapper nbt.Mapper) error {
			return mapper.MapString("", &level.ServerBrands[i])
		},
	))
	must(mapper.MapFloat("Data.SpawnAngle", &level.SpawnAngle))
	must(mapper.MapInt("Data.SpawnX", &level.Spawn.X))
	must(mapper.MapInt("Data.SpawnY", &level.Spawn.Y))
	must(mapper.MapInt("Data.SpawnZ", &level.Spawn.Z))
	must(mapper.MapLong("Data.Time", &level.Time))
	must(mapper.MapInt("Data.Version.Id", &level.Version.ID))
	must(mapper.MapString("Data.Version.Name", &level.Version.Name))
	must(mapper.MapByte("Data.Version.Snapshot", &level.Version.Snapshot))
	must(mapper.MapInt("Data.WanderingTraderSpawnChance", &level.WanderingTraderSpawnChance))
	must(mapper.MapInt("Data.WanderingTraderSpawnDelay", &level.WanderingTraderSpawnDelay))
	must(mapper.MapCustom("Data.WasModded", byteToBool(&level.WasModded)))
	must(mapper.MapCustom("Data.WorldGenSettings.bonus_chest", byteToBool(&level.WorldGenSettings.BonusChest)))
	must(mapper.MapCustom("Data.WorldGenSettings.dimensions", func(dimensions nbt.Tag) error {
		level.WorldGenSettings.DimensionsNBT = dimensions
		return nil
	}))
	must(mapper.MapCustom("Data.WorldGenSettings.generate_features", byteToBool(&level.WorldGenSettings.GenerateFeatures)))
	must(mapper.MapLong("Data.WorldGenSettings.seed", &level.WorldGenSettings.Seed))
	must(mapper.MapCustom("Data.allowCommands", byteToBool(&level.AllowCommands)))
	must(mapper.MapInt("Data.clearWeatherTime", &level.ClearWeatherTime))
	must(mapper.MapCustom("Data.hardcore", byteToBool(&level.Hardcore)))
	must(mapper.MapCustom("Data.initialized", byteToBool(&level.Initialized)))
	must(mapper.MapInt("Data.rainTime", &level.RainTime))
	must(mapper.MapCustom("Data.raining", byteToBool(&level.Raining)))
	must(mapper.MapInt("Data.thunderTime", &level.ThunderTime))
	must(mapper.MapCustom("Data.thundering", byteToBool(&level.Thundering)))
	return
}

func stringToBool(target *bool) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		val := tag.(*nbt.String).Value
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		*target = boolVal
		return nil
	}
}
func stringToInt(target *int) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		val := tag.(*nbt.String).Value
		intVal, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		*target = int(intVal)
		return nil
	}
}

func byteToBool(target *bool) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		*target = tag.(*nbt.Byte).Value != 0
		return nil
	}
}
