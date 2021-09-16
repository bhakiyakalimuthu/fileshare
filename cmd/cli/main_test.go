package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInputs(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "Should succeed if file name is valid",
			args:    []string{"file.png", "file.jpg", "file.gif"},
			want:    []string{"file.png", "file.jpg", "file.gif"},
			wantErr: false,
		},
		{
			name:    "Should succeed if file name is valid",
			args:    []string{"./file.png", "file.jpg", "file.gif"},
			want:    []string{"file.png", "file.jpg", "file.gif"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInputs(tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
