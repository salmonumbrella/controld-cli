package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "test", "account:test"},
		{"with dash", "my-account", "account:my-account"},
		{"with underscore", "my_account", "account:my_account"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, credentialKey(tt.input))
		})
	}
}

func TestParseCredentialKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOK   bool
	}{
		{"valid key", "account:myaccount", "myaccount", true},
		{"no prefix", "other:myaccount", "", false},
		{"just prefix", "account:", "", true},
		{"no colon", "accountmyaccount", "", false},
		{"wrong prefix", "other", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, ok := parseCredentialKey(tt.input)
			assert.Equal(t, tt.wantOK, ok)
			if tt.wantOK {
				assert.Equal(t, tt.wantName, name)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "test", "test"},
		{"uppercase", "TEST", "test"},
		{"mixed", "TeSt", "test"},
		{"with numbers", "Test123", "test123"},
		{"with leading space", "  test", "test"},
		{"with trailing space", "test  ", "test"},
		{"with both spaces", "  test  ", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalize(tt.input))
		})
	}
}
