package packs

import (
	"errors"
	"fmt"
	"slices"
)

var ErrInvalidSizes = errors.New("invalid sizes")

// Calculator is a package calculator service.
type Calculator interface {
	Calculate(int) map[int]int
}

type calculator struct {
	sizes []int
}

// NewCalculator validates the input sizes and creates a new calculator service.
func NewCalculator(sizes []int) (Calculator, error) {
	if err := validateSizes(sizes); err != nil {
		return nil, err
	}
	return &calculator{sizes: sizes}, nil
}

// Calculate calculates package sizes for a number of items.
func (c *calculator) Calculate(items int) map[int]int {
	res := make(map[int]int)
	for i := 0; i < len(c.sizes); i++ {
		size := c.sizes[len(c.sizes)-1-i]
		if items < size {
			continue
		}
		res[size] = items / size
		items = items % size
	}
	if items > 0 && len(c.sizes) > 0 {
		res[c.sizes[0]] += 1
	}
	return res
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
