// Package models contains the entities and the in-memory store
package models

import "sync"

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

// Get returns the product with the given ID.
func (s *Store) Get(id int) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	return p, ok
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

// Replace overwrites the whole product with the given ID.
func (s *Store) Replace(id int, p Product) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return Product{}, false
	}

	p.ID = id
	s.products[id] = p
	return p, true
}

// Update applies a partial change (PATCH) to an existing product.
func (s *Store) Update(id int, patch ProductPatch) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[id]
	if !ok {
		return Product{}, false
	}

	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Price != nil {
		p.Price = *patch.Price
	}
	if patch.Stock != nil {
		p.Stock = *patch.Stock
	}

	s.products[id] = p
	return p, true
}

// Delete removes the product with the given ID.
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return false
	}

	delete(s.products, id)
	return true
}
