// Code generated by "stringer -linecomment -output=serverbound_handshake_string.go -type=NextState"; DO NOT EDIT.

package packet

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NextStateStatus-1]
	_ = x[NextStateLogin-2]
}

const _NextState_name = "StatusLogin"

var _NextState_index = [...]uint8{0, 6, 11}

func (i NextState) String() string {
	i -= 1
	if i >= NextState(len(_NextState_index)-1) {
		return "NextState(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _NextState_name[_NextState_index[i]:_NextState_index[i+1]]
}
