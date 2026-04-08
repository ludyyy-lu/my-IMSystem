// Package session_test – presence store unit tests.
//
// These tests use an in-memory mock that satisfies the PresenceStore interface
// so they run without a real Redis instance.
package session_test

import (
	"context"
	"sync"
	"testing"

	"my-IMSystem/ws-gateway/internal/session"
)

// inMemPresenceStore is a thread-safe in-memory PresenceStore for tests.
type inMemPresenceStore struct {
	mu     sync.RWMutex
	online map[int64]struct{}
}

func newInMemPresenceStore() *inMemPresenceStore {
	return &inMemPresenceStore{online: make(map[int64]struct{})}
}

func (s *inMemPresenceStore) SetOnline(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.online[userID] = struct{}{}
	return nil
}

func (s *inMemPresenceStore) SetOffline(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.online, userID)
	return nil
}

func (s *inMemPresenceStore) IsOnline(_ context.Context, userID int64) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.online[userID]
	return ok, nil
}

// Compile-time assertion: *inMemPresenceStore satisfies session.PresenceStore.
var _ session.PresenceStore = (*inMemPresenceStore)(nil)

// TestPresenceStore_SetOnlineAndIsOnline verifies that SetOnline marks the user
// as online, and a subsequent IsOnline call returns true.
func TestPresenceStore_SetOnlineAndIsOnline(t *testing.T) {
	store := newInMemPresenceStore()
	ctx := context.Background()

	online, err := store.IsOnline(ctx, 1)
	if err != nil {
		t.Fatalf("IsOnline error: %v", err)
	}
	if online {
		t.Error("user should not be online before SetOnline")
	}

	if err := store.SetOnline(ctx, 1); err != nil {
		t.Fatalf("SetOnline error: %v", err)
	}

	online, err = store.IsOnline(ctx, 1)
	if err != nil {
		t.Fatalf("IsOnline error: %v", err)
	}
	if !online {
		t.Error("user should be online after SetOnline")
	}
}

// TestPresenceStore_SetOffline verifies that SetOffline clears the user's status.
func TestPresenceStore_SetOffline(t *testing.T) {
	store := newInMemPresenceStore()
	ctx := context.Background()

	_ = store.SetOnline(ctx, 2)
	if err := store.SetOffline(ctx, 2); err != nil {
		t.Fatalf("SetOffline error: %v", err)
	}

	online, err := store.IsOnline(ctx, 2)
	if err != nil {
		t.Fatalf("IsOnline error: %v", err)
	}
	if online {
		t.Error("user should be offline after SetOffline")
	}
}

// TestPresenceStore_SetOffline_Idempotent verifies that calling SetOffline when
// the user is already offline does not error.
func TestPresenceStore_SetOffline_Idempotent(t *testing.T) {
	store := newInMemPresenceStore()
	ctx := context.Background()

	// user was never online
	if err := store.SetOffline(ctx, 99); err != nil {
		t.Errorf("SetOffline on non-existent user should not error, got: %v", err)
	}
}

// TestPresenceStore_MultiUser verifies that presence is tracked independently
// per user.
func TestPresenceStore_MultiUser(t *testing.T) {
	store := newInMemPresenceStore()
	ctx := context.Background()

	_ = store.SetOnline(ctx, 10)
	_ = store.SetOnline(ctx, 20)
	_ = store.SetOffline(ctx, 10)

	if ok, _ := store.IsOnline(ctx, 10); ok {
		t.Error("user 10 should be offline")
	}
	if ok, _ := store.IsOnline(ctx, 20); !ok {
		t.Error("user 20 should still be online")
	}
}
