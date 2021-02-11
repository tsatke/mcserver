package block

import (
	"reflect"

	"github.com/tsatke/mcserver/game/id"
)

type Block interface {
	ID() id.ID
	Properties() map[string]Property
}

type BlockDescriptor struct {
	ID                  id.ID
	AvailableProperties []PropertyDescriptor
}

type Property interface {
	Name() string
	Value() interface{}
}

type PropertyDescriptor struct {
	Name          string
	Type          reflect.Type
	DefaultValue  interface{}
	AllowedValues []interface{}
}

type block struct {
	id         id.ID
	properties map[string]Property
}

func (b block) ID() id.ID                       { return b.id }
func (b block) Properties() map[string]Property { return b.properties }
