package wildcard

import (
	"encoding/json"
	"strings"

	"github.com/pucora/lura/v2/config"
)

// Namespace is the key used to look up plugin/wildcard config in ExtraConfig.
const Namespace = "plugin/wildcard"

// Config holds the wildcard plugin service-level configuration.
type Config struct {
	Prefix   string              `json:"prefix"`
	Mappings map[string][]string `json:"mappings"`
}

// ParseServiceConfig reads the plugin/wildcard config from a service-level ExtraConfig.
// It returns the parsed Config and true if the config was found and valid, or a zero-value
// Config and false otherwise.
func ParseServiceConfig(cfg config.ExtraConfig) (*Config, bool) {
	raw, ok := cfg[Namespace]
	if !ok {
		return nil, false
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, false
	}

	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, false
	}

	return &c, true
}

// IsWildcardEndpoint reports whether the given endpoint pattern is a wildcard.
// It returns true if endpoint is exactly "/*" or ends with "/*".
func IsWildcardEndpoint(endpoint string) bool {
	if endpoint == "/*" {
		return true
	}
	return strings.HasSuffix(endpoint, "/*")
}

// MatchWildcard looks up path in the provided mappings and returns the first
// matching mapped path and true. If no mapping matches, it returns an empty
// string and false.
//
// Matching is prefix-based: the request path is checked against each mapping
// key. If the path equals a key or starts with a key followed by "/", the
// first element of the associated slice is returned as the mapped path.
func MatchWildcard(path string, mappings map[string][]string) (string, bool) {
	for prefix, targets := range mappings {
		if len(targets) == 0 {
			continue
		}
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return targets[0], true
		}
	}
	return "", false
}
