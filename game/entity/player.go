package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

type (
	Player struct {
		Mob
		DataVersion            int
		PlayerGameType         int
		PreviousPlayerGameType int
		Score                  int
		Dimension              string
		SelectedItemSlot       int
		SelectedItem           InventoryItem
		SpawnDimension         string
		Spawn                  voxel.V3
		SpawnForced            bool
		SleepTimer             int16
		FoodLevel              int
		FoodExhaustionLevel    float32
		FoodSaturationLevel    float32
		FoodTickTimer          int
		XPLevel                int
		XPPercentage           float32
		XPTotal                int
		XPSeed                 int
		Inventory              []InventoryItem
		EnderItems             []InventoryItem
		Abilities              PlayerAbilities
		EnteredNetherPosition  [3]float64
		RootVehicle            PlayerRootVehicle
		ShoulderEntityLeft     Entity
		ShoulderEntityRight    Entity
		SeenCredits            bool
		RecipeBook             RecipeBook
	}

	RecipeBook struct {
		Recipes                             []id.ID
		ToBeDisplayed                       []id.ID
		IsFilteringCraftable                bool
		IsGuiOpen                           bool
		IsFurnaceFilteringCraftable         bool
		IsFurnaceGuiOpen                    bool
		IsBlastingFurnaceFilteringCraftable bool
		IsBlastingFurnaceGuiOpen            bool
		IsSmokerFilteringCraftable          bool
		IsSmokerGuiOpen                     bool
	}

	PlayerAbilities struct {
		WalkSpeed    float32
		FlySpeed     float32
		MayFly       bool
		Flying       bool
		Invulnerable bool
		MayBuild     bool
		InstaBuild   bool
	}

	PlayerRootVehicle struct {
		Attach uuid.UUID
		Entity Entity
	}
)

func PlayerFromNBTIntoPlayer(tag nbt.Tag, p *Player) (err error) {
	defer recoverAndSetErr(&err)

	mapper := nbt.NewSimpleMapper(tag)

	mob, err := decodeMob(mapper)
	if err != nil {
		return fmt.Errorf("decode mob: %w", err)
	}
	p.Mob = mob

	must(mapper.MapFloat("abilities.flySpeed", &p.Abilities.FlySpeed))
	must(mapper.MapCustom("abilities.flying", byteToBool(&p.Abilities.Flying)))
	must(mapper.MapCustom("abilities.instabuild", byteToBool(&p.Abilities.InstaBuild)))
	must(mapper.MapCustom("abilities.invulnerable", byteToBool(&p.Abilities.Invulnerable)))
	must(mapper.MapCustom("abilities.mayBuild", byteToBool(&p.Abilities.MayBuild)))
	must(mapper.MapCustom("abilities.mayfly", byteToBool(&p.Abilities.MayFly)))
	must(mapper.MapFloat("abilities.walkSpeed", &p.Abilities.WalkSpeed))
	must(mapper.MapInt("DataVersion", &p.DataVersion))
	must(mapper.MapString("Dimension", &p.Dimension))
	// TODO: must(mapper.Map[]Item("", &p.EnderItems))
	// TODO: must(mapper.Map[3]float64("", &p.EnteredNetherPosition))
	must(mapper.MapFloat("foodExhaustionLevel", &p.FoodExhaustionLevel))
	must(mapper.MapInt("foodLevel", &p.FoodLevel))
	must(mapper.MapFloat("foodSaturationLevel", &p.FoodSaturationLevel))
	must(mapper.MapInt("foodTickTimer", &p.FoodTickTimer))
	// TODO: must(mapper.Map[]Item("", &p.Inventory))
	must(mapper.MapInt("playerGameType", &p.PlayerGameType))
	must(mapper.MapInt("previousPlayerGameType", &p.PreviousPlayerGameType))
	// TODO: must(mapper.MapRecipeBook("", &p.RecipeBook))
	// TODO: must(mapper.MapPlayerRootVehicle("", &p.RootVehicle))
	must(mapper.MapInt("Score", &p.Score))
	must(mapper.MapCustom("seenCredits", byteToBool(&p.SeenCredits)))
	// TODO: must(mapper.MapItem("", &p.SelectedItem))
	must(mapper.MapInt("SelectedItemSlot", &p.SelectedItemSlot))
	// TODO: must(mapper.MapEntity("", &p.ShoulderEntityLeft))
	// TODO: must(mapper.MapEntity("", &p.ShoulderEntityRight))
	must(mapper.MapShort("SleepTimer", &p.SleepTimer))
	// TODO: must(mapper.Mapvoxel.V3("", &p.Spawn))
	_ = mapper.MapString("SpawnDimension", &p.SpawnDimension)
	_ = mapper.MapCustom("SpawnForced", byteToBool(&p.SpawnForced))
	must(mapper.MapInt("XpLevel", &p.XPLevel))
	_ = mapper.MapFloat("XpPercentage", &p.XPPercentage)
	must(mapper.MapInt("XpSeed", &p.XPSeed))
	must(mapper.MapInt("XpTotal", &p.XPTotal))
	return
}
