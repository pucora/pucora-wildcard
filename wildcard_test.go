package wildcard

import (
	"testing"

	"github.com/pucora/lura/v2/config"
)

func TestIsWildcardEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		want     bool
	}{
		{"/*", true},
		{"/api/*", true},
		{"/users/profile/*", true},
		{"/api", false},
		{"/users", false},
		{"", false},
		{"/api/*/extra", false},
		{"/", false},
		{"*", false},
	}

	for _, tc := range tests {
		got := IsWildcardEndpoint(tc.endpoint)
		if got != tc.want {
			t.Errorf("IsWildcardEndpoint(%q) = %v, want %v", tc.endpoint, got, tc.want)
		}
	}
}

func TestMatchWildcard(t *testing.T) {
	mappings := map[string][]string{
		"/api":   {"/v1/api", "/legacy/api"},
		"/users": {"/v1/users"},
	}

	tests := []struct {
		path     string
		wantPath string
		wantOK   bool
	}{
		{"/api", "/v1/api", true},
		{"/api/something", "/v1/api", true},
		{"/api/nested/path", "/v1/api", true},
		{"/users", "/v1/users", true},
		{"/users/123", "/v1/users", true},
		{"/other", "", false},
		{"/", "", false},
		{"", "", false},
		{"/apiv2", "", false},
	}

	for _, tc := range tests {
		gotPath, gotOK := MatchWildcard(tc.path, mappings)
		if gotOK != tc.wantOK || gotPath != tc.wantPath {
			t.Errorf("MatchWildcard(%q, mappings) = (%q, %v), want (%q, %v)",
				tc.path, gotPath, gotOK, tc.wantPath, tc.wantOK)
		}
	}
}

func TestMatchWildcard_EmptyMappings(t *testing.T) {
	path, ok := MatchWildcard("/api", nil)
	if ok || path != "" {
		t.Errorf("MatchWildcard with nil mappings: got (%q, %v), want (%q, false)", path, ok, "")
	}

	path, ok = MatchWildcard("/api", map[string][]string{})
	if ok || path != "" {
		t.Errorf("MatchWildcard with empty mappings: got (%q, %v), want (%q, false)", path, ok, "")
	}
}

func TestMatchWildcard_EmptyTargets(t *testing.T) {
	mappings := map[string][]string{
		"/api": {},
	}
	path, ok := MatchWildcard("/api", mappings)
	if ok || path != "" {
		t.Errorf("MatchWildcard with empty target slice: got (%q, %v), want (%q, false)", path, ok, "")
	}
}

func TestParseServiceConfig(t *testing.T) {
	t.Run("valid_config", func(t *testing.T) {
		extra := config.ExtraConfig{
			Namespace: map[string]interface{}{
				"prefix": "/__wildcard",
				"mappings": map[string]interface{}{
					"/api":   []interface{}{"/v1/api", "/legacy/api"},
					"/users": []interface{}{"/v1/users"},
				},
			},
		}
		cfg, ok := ParseServiceConfig(extra)
		if !ok {
			t.Fatal("expected ParseServiceConfig to return true")
		}
		if cfg.Prefix != "/__wildcard" {
			t.Errorf("Prefix = %q, want %q", cfg.Prefix, "/__wildcard")
		}
		if len(cfg.Mappings) != 2 {
			t.Errorf("Mappings length = %d, want 2", len(cfg.Mappings))
		}
		apiTargets := cfg.Mappings["/api"]
		if len(apiTargets) != 2 || apiTargets[0] != "/v1/api" {
			t.Errorf("Mappings[/api] = %v, want [/v1/api /legacy/api]", apiTargets)
		}
	})

	t.Run("missing_namespace", func(t *testing.T) {
		extra := config.ExtraConfig{}
		cfg, ok := ParseServiceConfig(extra)
		if ok || cfg != nil {
			t.Errorf("expected (nil, false) for missing namespace, got (%v, %v)", cfg, ok)
		}
	})
}
