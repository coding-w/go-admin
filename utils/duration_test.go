package utils

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name   string
		args   args
		target time.Duration
	}{
		{
			name:   "5h20m",
			args:   args{d: "5h20m"},
			target: time.Hour*5 + time.Minute*20,
		},
		{
			name:   "1d5h20m",
			args:   args{d: "1d5h20m"},
			target: time.Hour*24 + time.Hour*5 + time.Minute*20,
		},
		{
			name:   "1d",
			args:   args{"1d"},
			target: time.Hour * 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := ParseDuration(tt.args.d)
			if err != nil {
				t.Errorf("ParseDuration() error = %v", err)
				return
			}
			if duration != tt.target {
				t.Errorf("ParseDuration() = %v, want %v", duration, tt.target)
			}
		})
	}
}
