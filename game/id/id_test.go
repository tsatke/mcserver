package id

import (
	"reflect"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ID
	}{
		{
			"empty",
			"",
			ID{"", ""},
		},
		{
			"default namespace",
			"air",
			ID{"minecraft", "air"},
		},
		{
			"explicit namespace",
			"minecraft:grass",
			ID{"minecraft", "grass"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseID(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseID() = %v, want %v", got, tt.want)
			}
		})
	}
}
