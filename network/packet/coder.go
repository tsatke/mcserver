package packet

import (
	"encoding/binary"
	"fmt"
)

// Data type sizes in bytes.
const (
	BooleanSize       = 1
	ByteSize          = 1
	UnsignedByteSize  = 1
	ShortSize         = 2
	UnsignedShortSize = 2
	IntSize           = 4
	LongSize          = 8
	FloatSize         = 4
	DoubleSize        = 8
	VarIntMinSize     = 1
	VarIntMaxSize     = 5
	VarLongMinSize    = 1
	VarLongMaxSize    = 10
	PositionSize      = 8
	AngleSize         = 1
	UUIDSize          = 16
)

var (
	ByteOrder = binary.BigEndian
)

func recoverAndSetErr(err *error) {
	if rec := recover(); rec != nil {
		if recErr, ok := rec.(error); ok {
			*err = recErr
		} else {
			panic(rec)
		}
	}
}

func panicIffErr(fieldName string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", fieldName, err))
	}
}
