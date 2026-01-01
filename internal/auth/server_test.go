package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAccountName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid simple", "personal", false, ""},
		{"valid with dash", "my-account", false, ""},
		{"valid with underscore", "my_account", false, ""},
		{"valid with numbers", "account123", false, ""},
		{"valid mixed", "My_Account-123", false, ""},
		{"empty", "", true, "cannot be empty"},
		{"too long", strings.Repeat("a", 65), true, "too long"},
		{"max length", strings.Repeat("a", 64), false, ""},
		{"invalid space", "my account", true, "invalid characters"},
		{"invalid special", "my@account", true, "invalid characters"},
		{"invalid unicode", "账户", true, "invalid characters"},
		{"invalid dot", "my.account", true, "invalid characters"},
		{"invalid slash", "my/account", true, "invalid characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAPIToken(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid token", "api.abc123xyz", false, ""},
		{"valid long token", "api." + strings.Repeat("x", 100), false, ""},
		{"empty", "", true, "cannot be empty"},
		{"no prefix", "abc123xyz", true, "must start with 'api.'"},
		{"wrong prefix", "key.abc123", true, "must start with 'api.'"},
		{"just prefix", "api.", false, ""},
		{"too long", "api." + strings.Repeat("x", 253), true, "too long"},
		{"max length", "api." + strings.Repeat("x", 252), false, ""},
		{"partial prefix", "api", true, "must start with 'api.'"},
		{"similar prefix", "apix.token", true, "must start with 'api.'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIToken(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
