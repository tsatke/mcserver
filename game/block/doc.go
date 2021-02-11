// Package block provides means for registering blocks and creating blocks
// from registered block descriptors.
// A BlockDescriptor is a description of a block, kind of a blueprint.
// One can register a BlockDescriptor, which will allow for later creation
// of a block with the ID in the block descriptor.
//
//	id := id.ParseID("plugin:block")
//	myBlockDescriptor := BlockDescriptor{ID: id}
//	block.Must(block.RegisterBlock(myBlockDescriptor)) // panic on fail
//	block.Create(id)                                   // will succeed
//
// A block can not be created, if there's no block descriptor for it.
// Almost the same goes for properties that are supported for a block.
// A PropertyDescriptor doesn't have to be explicitly registered, but is
// registered with a BlockDescriptor. If one tries to create a block with
// a property whose PropertyDescriptor is not listed in the BlockDescriptor,
// the creation will fail.
//
// One can only register one BlockDescriptor per block ID and per numeric ID.
// Registering a duplicate block descriptor will result in an error.
package block
