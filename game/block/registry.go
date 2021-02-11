package block

import (
	"fmt"
	"reflect"

	"github.com/tsatke/mcserver/game/id"
)

var (
	blocksByID          = make(map[id.ID]BlockDescriptor)
	propertyDescriptors = make(map[string]PropertyDescriptor)
)

// Must panics if the given error is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func RegisterBlock(desc BlockDescriptor) error {
	if _, ok := blocksByID[desc.ID]; ok {
		return fmt.Errorf("block descriptor for block %s already exists", desc.ID)
	}
	blocksByID[desc.ID] = desc

	for _, propertyDesc := range desc.AvailableProperties {
		existing, ok := propertyDescriptors[propertyDesc.Name]
		if ok {
			// property descriptor already registered, check that they're equal
			if existing.Type != propertyDesc.Type {
				return fmt.Errorf("property descriptor %s already registered with different value type %s, while itself has value type %s", propertyDesc.Name, existing.Type.Name(), propertyDesc.Type.Name())
			}
			if !reflect.DeepEqual(existing.DefaultValue, propertyDesc.DefaultValue) {
				return fmt.Errorf("property descriptor %s already registered with different default value %v, while itself has default value %v", propertyDesc.Name, existing.DefaultValue, propertyDesc.DefaultValue)
			}
			// TODO: check AllowedValues
		}
		propertyDescriptors[propertyDesc.Name] = propertyDesc
	}
	return nil
}

func DescriptorForPropertyName(name string) (PropertyDescriptor, bool) {
	desc, ok := propertyDescriptors[name]
	return desc, ok
}

func Create(id id.ID, properties ...Property) (Block, error) {
	desc, ok := blocksByID[id]
	if !ok {
		return nil, fmt.Errorf("no block descriptor registered for id %s", id)
	}

	b := block{
		id:         desc.ID,
		properties: make(map[string]Property),
	}
	for _, property := range properties {
		// check if property is available for this block
		found := false
		for _, p := range desc.AvailableProperties {
			propertyDesc, ok := propertyDescriptors[property.Name()]
			if !ok {
				return nil, fmt.Errorf("no property descriptor for property %q", p.Name)
			}
			if p.Name == propertyDesc.Name {
				found = true
				break
			}
		}
		// set the property if available
		if !found {
			// property is not available
			return nil, fmt.Errorf("property %s is not available for block %s", property.Name(), desc.ID)
		}
		b.properties[property.Name()] = property
	}
	return b, nil
}

func CreateFromDescriptor(desc BlockDescriptor) (Block, error) {
	b, err := Create(desc.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to create block from block descriptor: %w", err)
	}
	return b, nil
}
