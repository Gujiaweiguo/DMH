package order

import (
	"context"
	"strings"
)

func hasVerificationPermission(ctx context.Context) bool {
	if ctx == nil {
		return true
	}

	rolesVal := ctx.Value("roles")
	if rolesVal == nil {
		// Keep backward compatibility for non-authenticated legacy flows.
		return true
	}

	roles := normalizeRoles(rolesVal)
	if len(roles) == 0 {
		return false
	}

	for _, role := range roles {
		if role == "platform_admin" || role == "brand_admin" {
			return true
		}
	}

	return false
}

func normalizeRoles(raw interface{}) []string {
	switch v := raw.(type) {
	case []string:
		return v
	case []interface{}:
		roles := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					roles = append(roles, trimmed)
				}
			}
		}
		return roles
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil
		}
		return []string{trimmed}
	default:
		return nil
	}
}
