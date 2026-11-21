// Package models contains the entities and the in-memory store
package models

import "sync"

// Store is an in-memory data store for products.
type Store struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

// NewStore creates an empty store ready to use.
func NewStore() *Store {
	return &Store{
		products: make(map[int]Product),
		nextID:   1,
	}
}

// List returns all products as a slice.
func (s *Store) List() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}
	return products
}

// Create stores a new product assigning it an auto-increment ID.
func (s *Store) Create(p Product) Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	s.nextID++
	s.products[p.ID] = p
	return p
}
