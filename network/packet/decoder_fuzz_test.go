package packet

import "bytes"

func FuzzDecodeHandshake(data []byte) error {
	return fuzzDecode(data, PhaseHandshaking)
}

func fuzzDecode(data []byte, phase Phase) error {
	p, err := Decode(bytes.NewReader(data), phase)
	if err != nil {
		return err
	}
	if v, ok := p.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
