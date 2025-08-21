package db

import "testing"

func TestClose(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Close"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Close()
		})
	}
}
