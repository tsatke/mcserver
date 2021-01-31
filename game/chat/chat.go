package chat

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
