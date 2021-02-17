// +build gofuzz

package packet

import "bytes"

func FuzzDecodeHandshake(data []byte) {
	fuzzDecode(data, PhaseHandshaking)
	fuzzDecode(data, PhaseStatus)
	fuzzDecode(data, PhaseLogin)
	fuzzDecode(data, PhasePlay)
}

func fuzzDecode(data []byte, phase Phase) {
	p, err := Decode(bytes.NewReader(data), phase)
	if err != nil {
		return
	}
	if v, ok := p.(Validator); ok {
		if err := v.Validate(); err != nil {
			return
		}
	}
	return nil
}
