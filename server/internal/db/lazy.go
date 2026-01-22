package db

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// LazyDatabase wraps Database to provide lazy connection with on-demand retry.
// It attempts to connect only when needed and maintains the connection once established.
type LazyDatabase struct {
	mu      sync.RWMutex
	db      *Database
	connErr error
	ctx     context.Context
}

// NewLazyDatabase creates a new LazyDatabase that will connect on first use.
func NewLazyDatabase(ctx context.Context) *LazyDatabase {
	return &LazyDatabase{
		ctx: ctx,
	}
}

// getOrConnect returns the database connection, attempting to connect if not already connected.
func (ld *LazyDatabase) getOrConnect() (*Database, error) {
	// Fast path: check if already connected with read lock
	ld.mu.RLock()
	if ld.db != nil {
		db := ld.db
		ld.mu.RUnlock()
		return db, nil
	}
	ld.mu.RUnlock()

	// Slow path: attempt connection with write lock
	ld.mu.Lock()
	defer ld.mu.Unlock()

	// Double-check after acquiring write lock
	if ld.db != nil {
		return ld.db, nil
	}

	// Attempt to connect
	db, err := NewDatabase(ld.ctx)
	if err != nil {
		ld.connErr = err
		log.Printf("Database connection attempt failed: %v", err)
		return nil, fmt.Errorf("database temporarily unavailable: %w", err)
	}

	ld.db = db
	ld.connErr = nil
	log.Println("Database connected successfully")
	return ld.db, nil
}

// Close closes the database connection if it exists.
func (ld *LazyDatabase) Close() {
	ld.mu.Lock()
	defer ld.mu.Unlock()
	if ld.db != nil {
		ld.db.Close()
		ld.db = nil
	}
}

// GetDatabase returns the underlying database connection, attempting to connect if necessary.
// This provides direct access to all repositories and methods.
func (ld *LazyDatabase) GetDatabase() (*Database, error) {
	return ld.getOrConnect()
}
