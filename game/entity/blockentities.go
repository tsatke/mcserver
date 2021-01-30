package entity

type (
	Banner struct {
		*BlockEntityData
		CustomName string
		Patterns   []BannerPattern
	}

	BannerPattern struct {
		Color   Color
		Pattern string
	}

	Barrel struct {
		*BlockEntityData
		CustomName    string
		Lock          string
		Items         []InventoryItem
		LootTable     interface{} // to be done
		LootTableSeed int64
	}

	// to be continued
)
