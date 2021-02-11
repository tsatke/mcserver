package loottable

import (
	"encoding/json"
	"io"
)

type Table struct {
	Type  string `json:"type"`
	Pools []Pool `json:"pools"`
}

type Count struct {
	Type string  `json:"type"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

type Function struct {
	Function string `json:"function"`
	Count    Count  `json:"count"`
	Add      bool   `json:"add,omitempty"`
}

type Entry struct {
	Type      string     `json:"type"`
	Functions []Function `json:"functions"`
	Name      string     `json:"name"`
}

type Condition struct {
	Condition         string  `json:"condition"`
	Chance            float64 `json:"chance,omitempty"`
	LootingMultiplier float64 `json:"looting_multiplier,omitempty"`
}

type Pool struct {
	Rolls      float64     `json:"rolls"`
	BonusRolls float64     `json:"bonus_rolls"`
	Entries    []Entry     `json:"entries"`
	Conditions []Condition `json:"conditions,omitempty"`
}

func FromJSONReader(rd io.Reader) (tbl Table, err error) {
	err = json.NewDecoder(rd).Decode(&tbl)
	return
}
