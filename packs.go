package packs

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

var ErrInvalidSizes = errors.New("invalid sizes")

// Calculator is a package calculator service.
type Calculator interface {
	Calculate(int) map[int]int
}

type calculator struct {
	sizeStore SizeStore
}

// NewCalculator validates the input sizes and creates a new calculator service.
func NewCalculator(sizeStore SizeStore) Calculator {
	return &calculator{sizeStore: sizeStore}
}

// Calculate calculates package sizes for a number of items.
func (c *calculator) Calculate(items int) map[int]int {
	sizes := c.sizeStore.Get()
	res := make(map[int]int)
	for i := 0; i < len(sizes); i++ {
		size := sizes[len(sizes)-1-i]
		if items < size {
			continue
		}
		res[size] = items / size
		items = items % size
	}
	if items > 0 && len(sizes) > 0 {
		res[sizes[0]] += 1
	}
	return res
}

// SizeStore stores sizes.
type SizeStore interface {
	Get() []int
	Set([]int) error
}

type sizeStore struct {
	sizes []int
	m     sync.RWMutex
}

// NewSizeStore validates the input sizes and creates a new SizeStore.
func NewSizeStore(sizes []int) (SizeStore, error) {
	s := &sizeStore{sizes: sizes}
	s.Set(sizes)
	return s, nil
}

// Get returns the sizes.
func (s *sizeStore) Get() []int {
	s.m.RLock()
	defer s.m.RUnlock()

	return append([]int{}, s.sizes...)
}

// Set validates and sets the sizes.
func (s *sizeStore) Set(sizes []int) error {
	s.m.Lock()
	defer s.m.Unlock()

	if err := validateSizes(sizes); err != nil {
		return err
	}
	s.sizes = append([]int{}, sizes...)
	return nil
}

func validateSizes(sizes []int) error {
	if len(sizes) == 0 {
		return errors.Join(ErrInvalidSizes, fmt.Errorf("must not be empty"))
	}
	if !slices.IsSorted(sizes) {
		return errors.Join(ErrInvalidSizes, fmt.Errorf("must be sorted"))
	}
	dup := make(map[int]struct{})
	for _, s := range sizes {
		if s <= 0 {
			return errors.Join(ErrInvalidSizes, fmt.Errorf("size must be greater than 0: %d", s))
		}
		if _, ok := dup[s]; ok {
			return errors.Join(ErrInvalidSizes, fmt.Errorf("must not contain duplicates: %d", s))
		}
		dup[s] = struct{}{}
	}
	return nil
}
