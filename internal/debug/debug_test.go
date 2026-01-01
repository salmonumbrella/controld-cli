package debug

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDebug(t *testing.T) {
	t.Run("sets debug true", func(t *testing.T) {
		ctx := WithDebug(context.Background(), true)
		assert.True(t, IsDebug(ctx))
	})

	t.Run("sets debug false", func(t *testing.T) {
		ctx := WithDebug(context.Background(), false)
		assert.False(t, IsDebug(ctx))
	})
}

func TestIsDebug(t *testing.T) {
	t.Run("default is false", func(t *testing.T) {
		assert.False(t, IsDebug(context.Background()))
	})
}

func TestSetupLogger(t *testing.T) {
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		SetupLogger(true)
	})
	assert.NotPanics(t, func() {
		SetupLogger(false)
	})
}
