package ui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		color string
	}{
		{"auto", "auto"},
		{"always", "always"},
		{"never", "never"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.color)
			assert.NotNil(t, u)
			assert.Equal(t, tt.color, u.color)
		})
	}
}

func TestWithUI(t *testing.T) {
	u := New("never")
	ctx := WithUI(context.Background(), u)

	got := FromContext(ctx)
	assert.Equal(t, u, got)
}

func TestFromContext(t *testing.T) {
	t.Run("returns UI from context", func(t *testing.T) {
		u := New("always")
		ctx := WithUI(context.Background(), u)
		assert.Equal(t, u, FromContext(ctx))
	})

	t.Run("returns default UI when not in context", func(t *testing.T) {
		got := FromContext(context.Background())
		assert.NotNil(t, got)
		assert.Equal(t, "auto", got.color)
	})
}

func TestUseColor(t *testing.T) {
	tests := []struct {
		name      string
		color     string
		wantColor bool
	}{
		{"never disables", "never", false},
		{"always enables", "always", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.color)
			assert.Equal(t, tt.wantColor, u.useColor())
		})
	}
}

func TestUseColorAuto(t *testing.T) {
	u := New("auto")
	// auto mode depends on terminal detection, just verify it doesn't panic
	_ = u.useColor()
}

func TestOutputMethods(t *testing.T) {
	// Output methods write to stderr, just verify they don't panic
	u := New("never")

	assert.NotPanics(t, func() {
		u.Success("test success")
	})
	assert.NotPanics(t, func() {
		u.Error("test error")
	})
	assert.NotPanics(t, func() {
		u.Info("test info")
	})
	assert.NotPanics(t, func() {
		u.Warn("test warning")
	})
}
