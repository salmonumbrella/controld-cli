package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveToken(t *testing.T) {
	t.Run("flag takes precedence", func(t *testing.T) {
		t.Setenv("CONTROLD_API_TOKEN", "env-token")

		cfg := ClientConfig{
			Token: "flag-token",
		}

		assert.Equal(t, "flag-token", resolveToken(cfg))
	})

	t.Run("env var used when no flag", func(t *testing.T) {
		t.Setenv("CONTROLD_API_TOKEN", "env-token")

		cfg := ClientConfig{}

		assert.Equal(t, "env-token", resolveToken(cfg))
	})

	t.Run("flag takes precedence over keyring", func(t *testing.T) {
		// Even if keyring has a token, flag should win
		cfg := ClientConfig{
			Token: "flag-token",
		}

		assert.Equal(t, "flag-token", resolveToken(cfg))
	})

	t.Run("env var takes precedence over keyring", func(t *testing.T) {
		t.Setenv("CONTROLD_API_TOKEN", "env-token")

		// Even if keyring has a token, env should win
		cfg := ClientConfig{}

		assert.Equal(t, "env-token", resolveToken(cfg))
	})
}

func TestNewClient(t *testing.T) {
	t.Run("creates client with flag token", func(t *testing.T) {
		cfg := ClientConfig{Token: "api.test123"}
		client, err := NewClient(context.Background(), cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("creates client with env token", func(t *testing.T) {
		t.Setenv("CONTROLD_API_TOKEN", "api.envtoken456")

		cfg := ClientConfig{}
		client, err := NewClient(context.Background(), cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("flag token takes precedence over env", func(t *testing.T) {
		t.Setenv("CONTROLD_API_TOKEN", "env-token")

		cfg := ClientConfig{Token: "flag-token"}
		client, err := NewClient(context.Background(), cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
