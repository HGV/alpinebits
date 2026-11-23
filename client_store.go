package alpinebits

import (
	"context"
	"sync"
)

// ClientStore keeps the result of each client's last handshake, so the router
// can hold its later requests to what was negotiated. The router identifies a
// client by the user of its basic access authentication; a store may instead
// key by anything ctx carries, such as values set by an authentication
// middleware.
type ClientStore interface {
	// Load returns what the client negotiated, by version; ok is false if it
	// never did a handshake.
	Load(ctx context.Context, user string) (negotiated map[string]ActionCapabilities, ok bool, err error)
	// Save replaces what the client negotiated.
	Save(ctx context.Context, user string, negotiated map[string]ActionCapabilities) error
}

// InMemoryClientStore is a ClientStore for tests and single instance servers:
// it is lost on restart and not shared between instances.
type InMemoryClientStore struct {
	mu   sync.RWMutex
	data map[string]map[string]ActionCapabilities
}

func NewInMemoryClientStore() *InMemoryClientStore {
	return &InMemoryClientStore{data: make(map[string]map[string]ActionCapabilities)}
}

func (s *InMemoryClientStore) Load(_ context.Context, user string) (map[string]ActionCapabilities, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	negotiated, ok := s.data[user]
	return negotiated, ok, nil
}

func (s *InMemoryClientStore) Save(_ context.Context, user string, negotiated map[string]ActionCapabilities) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[user] = negotiated
	return nil
}
