package entity

//go:generate stringer -linecomment -output effect_id_string.go -type=EffectID

type EffectID byte

const (
	EffectIDInvalid          EffectID = iota // invalid
	EffectIDSpeed                            // speed
	EffectIDSlowness                         // slowness
	EffectIDHaste                            // haste
	EffectIDMiningFatigue                    // mining_fatigue
	EffectIDStrength                         // strength
	EffectIDInstantHealth                    // instant_health
	EffectIDInstantDamage                    // instant_damage
	EffectIDJumpBoost                        // jump_boost
	EffectIDNausea                           // nausea
	EffectIDRegeneration                     // regeneration
	EffectIDResistance                       // resistance
	EffectIDFireResistance                   // fire_resistance
	EffectIDWaterBreathing                   // water_breathing
	EffectIDInvisibility                     // invisibility
	EffectIDBlindness                        // blindness
	EffectIDNightVision                      // night_vision
	EffectIDHunger                           // hunger
	EffectIDWeakness                         // weakness
	EffectIDPoison                           // poison
	EffectIDWither                           // wither
	EffectIDHealthBoost                      // health_boost
	EffectIDAbsorption                       // absorption
	EffectIDSaturation                       // saturation
	EffectIDGlowing                          // glowing
	EffectIDLevitation                       // levitation
	EffectIDLuck                             // luck
	EffectIDBadLuck                          // unluck
	EffectIDSlowFalling                      // slow_falling
	EffectIDConduitPower                     // conduit_power
	EffectIDDolphinsGrace                    // dolphins_grace
	EffectIDBadOmen                          // bad_omen
	EffectIDHeroOfTheVillage                 // hero_of_the_village
)
