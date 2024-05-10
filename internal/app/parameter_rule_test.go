package app

import (
	"fmt"
	"testing"
)

func TestParameterRuleIsCovers(t *testing.T) {
	tests := []struct {
		r1   ParameterRule
		r2   ParameterRule
		want bool
	}{
		{
			r1: ParameterRule{
				Path:  "/foo/v1",
				Level: ParameterLevelStrict,
			},
			r2: ParameterRule{
				Path:  "/foo/v1",
				Level: ParameterLevelStrict,
			},
			want: true,
		},
		{
			r1: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelUnder,
			},
			r2: ParameterRule{
				Path:  "/foo/v1",
				Level: ParameterLevelStrict,
			},
			want: true,
		},
		{
			r1: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelUnder,
			},
			r2: ParameterRule{
				Path:  "/foo/v1/value",
				Level: ParameterLevelStrict,
			},
			want: false,
		},
		{
			r1: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelAll,
			},
			r2: ParameterRule{
				Path:  "/foo/v1/value",
				Level: ParameterLevelStrict,
			},
			want: true,
		},
		{
			r1: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelAll,
			},
			r2: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelUnder,
			},
			want: true,
		},
		{
			r1: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelUnder,
			},
			r2: ParameterRule{
				Path:  "/foo/",
				Level: ParameterLevelAll,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s covers %s -> %t", tt.r1, tt.r2, tt.want), func(t *testing.T) {
			if tt.r1.IsCovers(tt.r2) != tt.want {
				t.Errorf("unexpected result")
			}
		})
	}
}
