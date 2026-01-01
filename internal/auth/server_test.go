package auth

import (
	"strings"
	"testing"
	"time"

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

func TestRateLimiter(t *testing.T) {
	t.Run("allows requests under limit", func(t *testing.T) {
		rl := newRateLimiter(3, time.Minute)

		for i := 0; i < 3; i++ {
			err := rl.check("127.0.0.1", "/test")
			assert.NoError(t, err, "request %d should be allowed", i+1)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		rl := newRateLimiter(2, time.Minute)

		assert.NoError(t, rl.check("127.0.0.1", "/test"))
		assert.NoError(t, rl.check("127.0.0.1", "/test"))
		err := rl.check("127.0.0.1", "/test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many attempts")
	})

	t.Run("tracks different endpoints separately", func(t *testing.T) {
		rl := newRateLimiter(1, time.Minute)

		assert.NoError(t, rl.check("127.0.0.1", "/a"))
		assert.NoError(t, rl.check("127.0.0.1", "/b"))
		assert.Error(t, rl.check("127.0.0.1", "/a"))
		assert.Error(t, rl.check("127.0.0.1", "/b"))
	})

	t.Run("tracks different IPs separately", func(t *testing.T) {
		rl := newRateLimiter(1, time.Minute)

		assert.NoError(t, rl.check("127.0.0.1", "/test"))
		assert.NoError(t, rl.check("127.0.0.2", "/test"))
		assert.Error(t, rl.check("127.0.0.1", "/test"))
	})

	t.Run("resets after window expires", func(t *testing.T) {
		rl := newRateLimiter(1, 10*time.Millisecond)

		assert.NoError(t, rl.check("127.0.0.1", "/test"))
		assert.Error(t, rl.check("127.0.0.1", "/test"))

		time.Sleep(15 * time.Millisecond)

		assert.NoError(t, rl.check("127.0.0.1", "/test"))
	})
}

func TestRateLimiterCleanup(t *testing.T) {
	rl := newRateLimiter(1, 10*time.Millisecond)

	_ = rl.check("127.0.0.1", "/test")
	assert.Len(t, rl.attempts, 1)

	time.Sleep(15 * time.Millisecond)
	rl.cleanup()

	assert.Len(t, rl.attempts, 0)
}
