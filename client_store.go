package alpinebits

import "sync"

// ClientStore stores negotiated capabilities per client and version.
type ClientStore interface {
	Get(clientID, version string) ActionCapabilities
	Set(clientID, version string, caps ActionCapabilities)
}

// clientKey creates a composite key for clientID + version.
func clientKey(clientID, version string) string {
	return clientID + ":" + version
}

// InMemoryClientStore is a thread-safe in-memory implementation of ClientStore.
type InMemoryClientStore struct {
	mu   sync.RWMutex
	data map[string]ActionCapabilities
}

// NewInMemoryClientStore creates a new in-memory client store.
func NewInMemoryClientStore() *InMemoryClientStore {
	return &InMemoryClientStore{
		data: make(map[string]ActionCapabilities),
	}
}

// Get retrieves negotiated capabilities for a client and version.
func (s *InMemoryClientStore) Get(clientID, version string) ActionCapabilities {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[clientKey(clientID, version)]
}

// Set stores negotiated capabilities for a client and version.
func (s *InMemoryClientStore) Set(clientID, version string, caps ActionCapabilities) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[clientKey(clientID, version)] = caps
}
