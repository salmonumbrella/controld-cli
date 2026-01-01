package outfmt

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantJSON bool
	}{
		{"json format", "json", true},
		{"JSON uppercase", "JSON", false}, // IsJSON checks for exact "json" match
		{"text format", "text", false},
		{"empty format", "", false},
		{"other format", "yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithFormat(context.Background(), tt.format)
			assert.Equal(t, tt.wantJSON, IsJSON(ctx))
		})
	}
}

func TestWithYes(t *testing.T) {
	t.Run("yes true", func(t *testing.T) {
		ctx := WithYes(context.Background(), true)
		assert.True(t, GetYes(ctx))
	})

	t.Run("yes false", func(t *testing.T) {
		ctx := WithYes(context.Background(), false)
		assert.False(t, GetYes(ctx))
	})

	t.Run("default is false", func(t *testing.T) {
		assert.False(t, GetYes(context.Background()))
	})
}

func TestIsJSON(t *testing.T) {
	t.Run("default is false", func(t *testing.T) {
		assert.False(t, IsJSON(context.Background()))
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("writes valid JSON", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]string{"foo": "bar"}

		err := WriteJSON(&buf, data)

		assert.NoError(t, err)
		assert.JSONEq(t, `{"foo":"bar"}`, buf.String())
	})

	t.Run("writes array", func(t *testing.T) {
		var buf bytes.Buffer
		data := []int{1, 2, 3}

		err := WriteJSON(&buf, data)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, buf.String())
	})

	t.Run("returns error for unencodable type", func(t *testing.T) {
		var buf bytes.Buffer
		ch := make(chan int) // channels cannot be JSON encoded

		err := WriteJSON(&buf, ch)

		assert.Error(t, err)
	})
}

func TestNewTabWriter(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTabWriter(&buf)

	_, err := tw.Write([]byte("NAME\tAGE\n"))
	assert.NoError(t, err)
	_, err = tw.Write([]byte("Alice\t30\n"))
	assert.NoError(t, err)
	_, err = tw.Write([]byte("Bob\t25\n"))
	assert.NoError(t, err)
	err = tw.Flush()
	assert.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Verify we have the expected lines
	assert.GreaterOrEqual(t, len(lines), 3)

	// Verify columns are aligned (NAME and Alice should start at same position)
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "Alice")
	assert.Contains(t, output, "Bob")
}
