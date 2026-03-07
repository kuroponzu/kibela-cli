package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		team      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid config",
			token:   "test-token",
			team:    "test-team",
			wantErr: false,
		},
		{
			name:      "missing token",
			token:     "",
			team:      "test-team",
			wantErr:   true,
			errSubstr: "KIBELA_TOKEN",
		},
		{
			name:      "missing team",
			token:     "test-token",
			team:      "",
			wantErr:   true,
			errSubstr: "KIBELA_TEAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear existing env vars
			os.Unsetenv("KIBELA_TOKEN")
			os.Unsetenv("KIBELA_TEAM")

			// Set test env vars
			if tt.token != "" {
				os.Setenv("KIBELA_TOKEN", tt.token)
			}
			if tt.team != "" {
				os.Setenv("KIBELA_TEAM", tt.team)
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() error = nil, want error containing %q", tt.errSubstr)
					return
				}
				if tt.errSubstr != "" && !containsString(err.Error(), tt.errSubstr) {
					t.Errorf("Load() error = %q, want error containing %q", err.Error(), tt.errSubstr)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v", err)
				return
			}

			if cfg.Token != tt.token {
				t.Errorf("Load() Token = %q, want %q", cfg.Token, tt.token)
			}
			if cfg.Team != tt.team {
				t.Errorf("Load() Team = %q, want %q", cfg.Team, tt.team)
			}
		})
	}
}

func TestConfig_Endpoint(t *testing.T) {
	cfg := &Config{
		Token: "test-token",
		Team:  "my-team",
	}

	got := cfg.Endpoint()
	want := "https://my-team.kibe.la/api/v1"

	if got != want {
		t.Errorf("Endpoint() = %q, want %q", got, want)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
