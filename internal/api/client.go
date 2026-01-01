package api

import (
	"context"
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"

	"github.com/salmonumbrella/controld-cli/internal/config"
	"github.com/salmonumbrella/controld-cli/internal/debug"
	"github.com/salmonumbrella/controld-cli/internal/secrets"
)

type ClientConfig struct {
	Token   string
	Account string
}

func NewClient(ctx context.Context, cfg ClientConfig) (*controld.API, error) {
	token := resolveToken(cfg)
	if token == "" {
		return nil, fmt.Errorf("no API token found. Set %s or run: controld auth login", config.EnvToken)
	}

	opts := []controld.Option{}
	if debug.IsDebug(ctx) {
		opts = append(opts, controld.Debug(true))
	}

	return controld.New(token, opts...)
}

func resolveToken(cfg ClientConfig) string {
	// 1. Explicit token flag
	if cfg.Token != "" {
		return cfg.Token
	}

	// 2. Environment variable
	if token := os.Getenv(config.EnvToken); token != "" {
		return token
	}

	// 3. Keyring
	store, err := secrets.OpenDefault()
	if err != nil {
		return ""
	}

	// If account specified, use it
	if cfg.Account != "" {
		creds, err := store.Get(cfg.Account)
		if err != nil {
			return ""
		}
		return creds.Token
	}

	// Auto-select if only one account
	creds, err := store.List()
	if err != nil || len(creds) != 1 {
		return ""
	}
	return creds[0].Token
}
