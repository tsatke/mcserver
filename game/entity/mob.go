package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

type (
	Mob struct {
		Data

		Health              float32
		AbsorptionAmount    float32
		HurtTime            int16
		HurtByTimestamp     int
		DeathTime           int16
		FallFlying          bool
		SleepingX           int
		SleepingY           int
		SleepingZ           int
		TicksFrozen         int
		Brain               Brain
		Attributes          []Attribute
		ActiveEffects       []PotionEffect
		HandItems           []InventoryItem
		ArmorItems          []InventoryItem
		HandDropChanges     []float32
		ArmorDropChanges    []float32
		DeathLootTable      interface{} // to be done
		DeathLootTableSeed  int64
		CanPickUpLoot       bool
		NoAI                bool
		PersistenceRequired bool
		LeftHanded          bool
		Team                interface{} // to be done
		Leash               EntityLeash
	}

	Brain struct {
		Memories map[id.ID]interface{}
	}

	Attribute struct {
		Name      id.ID
		Base      float64
		Modifiers []Modifier
	}

	Modifier struct {
		Name      id.ID
		Amount    float64
		Operation int
		UUID      uuid.UUID
	}

	PotionEffect struct {
		EffectID      EffectID
		Amplifier     byte
		Duration      int
		Ambient       bool
		ShowParticles bool
		ShowIcon      bool
	}

	EntityLeash struct {
		HasUUID bool
		UUID    uuid.UUID
		X       int
		Y       int
		Z       int
	}
)

// From https://minecraft.gamepedia.com/Chunk_format#Entity_format
type (
	Bat struct {
		Mob
		BatFlags int8
	}

	Bee struct {
		Mob
		CanBreed
		CanBeAngry
		HivePos                    voxel.V3
		FlowerPos                  voxel.V3
		HasNectar                  bool
		HasStung                   bool
		TicksSincePollination      int
		CannotEnterHiveTicks       int
		CropsGrownSincePollination int
	}

	Blaze struct {
		Mob
	}

	Cat struct {
		Mob
		CatType CatType
	}

	CaveSpider struct {
		Mob
	}

	Chicken struct {
		Mob
		CanBreed
		IsChickenJockey bool
		EggLayTime      int
	}

	Cod struct {
		Mob
		FromBucket bool
	}

	Cow struct {
		Mob
		CanBreed
	}

	Creeper struct {
		Mob
		Powered         bool
		ExplisionRadius int8
		Fuse            int16
		Ignited         bool
	}

	Pig struct {
		Mob
		CanBreed
		Saddle bool
	}

	Sheep struct {
		Mob
		CanBreed
		Sheared bool
		Color   int8
	}

	Skeleton struct {
		Mob
	}

	Wolf struct {
		Mob
		CanBreed
		CanBeTamed
		CanBeAngry
		CollarColor int8
	}

	// to be continued
)

func FromNBT(id id.ID, tag nbt.Tag) (e Entity, err error) {
	defer recoverAndSetErr(&err)

	compound, ok := tag.(*nbt.Compound)
	if !ok {
		return nil, fmt.Errorf("tag is not a compound (got %s)", tag.ID())
	}
	decoder, ok := nbtDecoders[id]
	if !ok {
		return nil, fmt.Errorf("no NBT decoder for id %s", id)
	}
	return decoder(nbt.NewSimpleMapper(compound))
}

func decodeData(mapper nbt.Mapper) (data Data, err error) {
	defer recoverAndSetErr(&err)

	// optional for players
	_ = mapper.MapCustom("id", func(tag nbt.Tag) error {
		data.ID = id.ParseID(tag.(*nbt.String).Value)
		return nil
	})
	must(mapper.MapList("Pos", func(int) {}, func(i int, mapper nbt.Mapper) error {
		return mapper.MapDouble("", &data.Pos[i])
	}))
	must(mapper.MapList("Motion", func(int) {}, func(i int, mapper nbt.Mapper) error {
		return mapper.MapDouble("", &data.Motion[i])
	}))
	must(mapper.MapList("Rotation", func(int) {}, func(i int, mapper nbt.Mapper) error {
		return mapper.MapFloat("", &data.Rotation[i])
	}))
	must(mapper.MapFloat("FallDistance", &data.FallDistance))
	must(mapper.MapShort("Fire", &data.Fire))
	must(mapper.MapShort("Air", &data.Air))
	must(mapper.MapCustom("OnGround", byteToBool(&data.OnGround)))
	_ = mapper.MapCustom("NoGravity", byteToBool(&data.NoGravity))
	must(mapper.MapCustom("Invulnerable", byteToBool(&data.Invulnerable)))
	must(mapper.MapInt("PortalCooldown", &data.PortalCooldown))
	must(mapper.MapCustom("UUID", intsToUUID(&data.UUID)))
	_ = mapper.MapString("CustomName", &data.CustomName)
	_ = mapper.MapCustom("CustomNameVisible", byteToBool(&data.CustomNameVisible))
	_ = mapper.MapCustom("Silent", byteToBool(&data.Silent))
	// TODO: passengers
	_ = mapper.MapCustom("Glowing", byteToBool(&data.Glowing))
	// TODO: tags
	return
}

func decodeMob(mapper nbt.Mapper) (mob Mob, err error) {
	defer recoverAndSetErr(&err)

	data, err := decodeData(mapper)
	if err != nil {
		return Mob{}, fmt.Errorf("decode data: %w", err)
	}
	mob.Data = data

	must(mapper.MapFloat("Health", &mob.Health))
	must(mapper.MapFloat("AbsorptionAmount", &mob.AbsorptionAmount))
	must(mapper.MapShort("HurtTime", &mob.HurtTime))
	must(mapper.MapInt("HurtByTimestamp", &mob.HurtByTimestamp))
	must(mapper.MapShort("DeathTime", &mob.DeathTime))
	must(mapper.MapCustom("FallFlying", byteToBool(&mob.FallFlying)))
	_ = mapper.MapInt("SleepingX", &mob.SleepingX)
	_ = mapper.MapInt("SleepingY", &mob.SleepingY)
	_ = mapper.MapInt("SleepingZ", &mob.SleepingZ)
	_ = mapper.MapInt("TickFrozen", &mob.TicksFrozen)
	// TODO: Brain
	// TODO: Attributes
	// TODO: ActiveEffects
	// TODO: HandItems
	// TODO: ArmorItems
	// TODO: HandDropChances
	// TODO: ArmorDropChances
	// TODO: DeathLootTable
	_ = mapper.MapLong("DeathLootTableSeed", &mob.DeathLootTableSeed)
	_ = mapper.MapCustom("CanPickUpLoot", byteToBool(&mob.CanPickUpLoot)) // optional for player
	_ = mapper.MapCustom("NoAI", byteToBool(&mob.NoAI))
	_ = mapper.MapCustom("PersistenceRequired", byteToBool(&mob.PersistenceRequired)) // optional for player
	_ = mapper.MapCustom("LeftHanded", byteToBool(&mob.LeftHanded))                   // optional for player
	// TODO: Team
	// TODO: Leash
	return
}

func decodeCanBreed(mapper nbt.Mapper) (cb CanBreed, err error) {
	defer recoverAndSetErr(&err)

	must(mapper.MapInt("InLove", &cb.InLove))
	must(mapper.MapInt("Age", &cb.Age))
	must(mapper.MapInt("ForcedAge", &cb.ForcedAge))
	_ = mapper.MapCustom("LoveCause", intsToUUID(&cb.LoveCause))

	return
}

func decodeCanBeAngry(mapper nbt.Mapper) (cba CanBeAngry, err error) {
	defer recoverAndSetErr(&err)

	must(mapper.MapInt("AngerTime", &cba.AngerTime))
	_ = mapper.MapCustom("AngryAt", intsToUUID(&cba.AngryAt))

	return
}

func decodeCanBeTamed(mapper nbt.Mapper) (cba CanBeTamed, err error) {
	defer recoverAndSetErr(&err)

	_ = mapper.MapCustom("Owner", intsToUUID(&cba.Owner))
	must(mapper.MapCustom("Sitting", byteToBool(&cba.Sitting)))

	return
}

func decodeBat(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	bat := Bat{
		Mob: mob,
	}
	must(mapper.MapByte("BatFlags", &bat.BatFlags))

	return &bat, nil
}

func decodeCreeper(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	creeper := Creeper{
		Mob: mob,
	}
	_ = mapper.MapCustom("powered", byteToBool(&creeper.Powered))
	must(mapper.MapByte("ExplosionRadius", &creeper.ExplisionRadius))
	must(mapper.MapShort("Fuse", &creeper.Fuse))
	must(mapper.MapCustom("ignited", byteToBool(&creeper.Ignited)))

	return &creeper, nil
}

func decodeCow(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	cb, err := decodeCanBreed(mapper)
	if err != nil {
		return nil, err
	}

	return &Cow{
		Mob:      mob,
		CanBreed: cb,
	}, nil
}

func decodePig(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	cb, err := decodeCanBreed(mapper)
	if err != nil {
		return nil, err
	}

	pig := Pig{
		Mob:      mob,
		CanBreed: cb,
	}

	must(mapper.MapCustom("Saddle", byteToBool(&pig.Saddle)))

	return &pig, nil
}

func decodeSheep(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	cb, err := decodeCanBreed(mapper)
	if err != nil {
		return nil, err
	}

	sheep := Sheep{
		Mob:      mob,
		CanBreed: cb,
	}

	must(mapper.MapCustom("Sheared", byteToBool(&sheep.Sheared)))
	must(mapper.MapByte("Color", &sheep.Color))

	return &sheep, nil
}

func decodeSkeleton(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}

	skeleton := Skeleton{
		Mob: mob,
	}
	return &skeleton, nil
}

func decodeWolf(mapper nbt.Mapper) (Entity, error) {
	mob, err := decodeMob(mapper)
	if err != nil {
		return nil, err
	}
	canBeAngry, err := decodeCanBeAngry(mapper)
	if err != nil {
		return nil, err
	}
	canBeTamed, err := decodeCanBeTamed(mapper)
	if err != nil {
		return nil, err
	}

	wolf := Wolf{
		Mob:        mob,
		CanBeAngry: canBeAngry,
		CanBeTamed: canBeTamed,
	}

	must(mapper.MapByte("CollarColor", &wolf.CollarColor))

	return &wolf, nil
}
