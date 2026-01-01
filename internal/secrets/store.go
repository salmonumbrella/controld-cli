package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/keyring"

	"github.com/salmonumbrella/controld-cli/internal/config"
)

type Store interface {
	Set(name string, token string) error
	Get(name string) (Credentials, error)
	Delete(name string) error
	List() ([]Credentials, error)
	Keys() ([]string, error)
}

type Credentials struct {
	Name      string    `json:"name"`
	Token     string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type storedCredentials struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type KeyringStore struct {
	ring keyring.Keyring
}

func OpenDefault() (Store, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: config.AppName,
	})
	if err != nil {
		return nil, err
	}
	return &KeyringStore{ring: ring}, nil
}

func (s *KeyringStore) Keys() ([]string, error) {
	return s.ring.Keys()
}

func (s *KeyringStore) Set(name string, token string) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	if token == "" {
		return fmt.Errorf("missing token")
	}

	payload, err := json.Marshal(storedCredentials{
		Token:     token,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return s.ring.Set(keyring.Item{
		Key:  credentialKey(name),
		Data: payload,
	})
}

func (s *KeyringStore) Get(name string) (Credentials, error) {
	name = normalize(name)
	if name == "" {
		return Credentials{}, fmt.Errorf("missing account name")
	}

	item, err := s.ring.Get(credentialKey(name))
	if err != nil {
		return Credentials{}, err
	}

	var stored storedCredentials
	if err := json.Unmarshal(item.Data, &stored); err != nil {
		return Credentials{}, err
	}

	return Credentials{
		Name:      name,
		Token:     stored.Token,
		CreatedAt: stored.CreatedAt,
	}, nil
}

func (s *KeyringStore) Delete(name string) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	return s.ring.Remove(credentialKey(name))
}

func (s *KeyringStore) List() ([]Credentials, error) {
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}

	var out []Credentials
	for _, k := range keys {
		name, ok := parseCredentialKey(k)
		if !ok {
			continue
		}
		creds, err := s.Get(name)
		if err != nil {
			continue
		}
		out = append(out, creds)
	}
	return out, nil
}

func credentialKey(name string) string {
	return fmt.Sprintf("account:%s", name)
}

func parseCredentialKey(k string) (string, bool) {
	const prefix = "account:"
	if !strings.HasPrefix(k, prefix) {
		return "", false
	}
	return strings.TrimPrefix(k, prefix), true
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
