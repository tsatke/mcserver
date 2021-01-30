package types

import (
	"encoding/json"
	"fmt"
	"io"
)

type (
	Chat struct {
		ChatFragment
		Extra []ChatFragment `json:"extra,omitempty"`
	}
	ChatFragment struct {
		Text          string `json:"text"`
		Bold          bool   `json:"bold,omitempty"`
		Italic        bool   `json:"italic,omitempty"`
		Underlined    bool   `json:"underlined,omitempty"`
		Strikethrough bool   `json:"strikethrough,omitempty"`
		Obfuscated    bool   `json:"obfuscated,omitempty"`
		Color         string `json:"color,omitempty"`
		Insertion     string `json:"insertion,omitempty"`
	}
)

func (ch *Chat) DecodeFrom(rd io.Reader) error {
	strVal := NewString("")
	if err := strVal.DecodeFrom(rd); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(*strVal), ch); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func (ch Chat) EncodeInto(w io.Writer) error {
	data, err := json.Marshal(ch)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return NewString(string(data)).EncodeInto(w)
}
