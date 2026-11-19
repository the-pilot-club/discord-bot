package tpc

import (
	"testing"

	"github.com/the-pilot-club/tpcgo"
)

func TestResolveCoreAPIEnvironment(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  tpcgo.Environment
	}{
		{name: "beta", value: "beta", want: tpcgo.EnvBeta},
		{name: "beta ignores case and whitespace", value: " Beta ", want: tpcgo.EnvBeta},
		{name: "production", value: "production", want: tpcgo.EnvProduction},
		{name: "empty defaults to production", value: "", want: tpcgo.EnvProduction},
		{name: "unknown defaults to production", value: "staging", want: tpcgo.EnvProduction},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveCoreAPIEnvironment(tt.value); got != tt.want {
				t.Fatalf("resolveCoreAPIEnvironment(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
