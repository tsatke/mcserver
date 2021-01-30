package block

import "reflect"

//go:generate stringer -type=Property

type PropertyDescriptor struct {
	Name          string
	Type          reflect.Type
	DefaultValue  interface{}
	AllowedValues []interface{}
}

func createPropertyDescriptor(name string, defaultValue interface{}, allowedValues ...interface{}) PropertyDescriptor {
	return PropertyDescriptor{
		Name:          name,
		Type:          reflect.TypeOf(defaultValue),
		DefaultValue:  defaultValue,
		AllowedValues: allowedValues,
	}
}

var (
	compassDirections = []interface{}{"east", "north", "south", "west"}
	AnvilFacing       = createPropertyDescriptor("facing", "north", compassDirections...)
	BambooAge         = createPropertyDescriptor("age", 0, 0, 1)
	BambooLeaves      = createPropertyDescriptor("leaves", "none", "large", "none", "small")
	// to be continued
)
