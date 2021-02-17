// +build gofuzz

package packet

import "bytes"

func Fuzz(data []byte) int {
	res := 0
	res += fuzzDecode(data, PhaseHandshaking)
	res += fuzzDecode(data, PhaseStatus)
	res += fuzzDecode(data, PhaseLogin)
	res += fuzzDecode(data, PhasePlay)
	if res > 0 {
		return 1
	}
	return 0
}

func fuzzDecode(data []byte, phase Phase) int {
	p, err := Decode(bytes.NewReader(data), phase)
	if err != nil {
		return 0
	}
	if v, ok := p.(Validator); ok {
		if err := v.Validate(); err != nil {
			return 0
		}
	}
	return 1
}
