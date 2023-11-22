package packs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/szkrstf/packs/mock"
)

func TestCalculate(t *testing.T) {
	s := calculator{sizeStore: &mock.SizeStore{
		GetFn: func() []int { return []int{250, 500, 1000, 2000, 5000} },
	}}

	tt := []struct {
		items int
		packs map[int]int
	}{
		{items: 0, packs: map[int]int{}},
		{items: 1, packs: map[int]int{250: 1}},
		{items: 250, packs: map[int]int{250: 1}},
		{items: 251, packs: map[int]int{250: 2}},
		{items: 501, packs: map[int]int{250: 1, 500: 1}},
		{items: 12001, packs: map[int]int{250: 1, 2000: 1, 5000: 2}},
	}

	for _, tc := range tt {
		if got, want := s.Calculate(tc.items), tc.packs; !reflect.DeepEqual(got, want) {
			t.Errorf("%v: got: %v; want: %v", tc.items, got, want)
		}
	}
}

func TestSizeStore(t *testing.T) {
	sizes := []int{1, 2}
	s := sizeStore{sizes: sizes}

	tt := []struct {
		input []int
		err   error
	}{
		{input: nil, err: ErrInvalidSizes},
		{input: []int{3, 2, 1}, err: ErrInvalidSizes},
		{input: []int{-1, 2, 3}, err: ErrInvalidSizes},
		{input: []int{1, 1, 1}, err: ErrInvalidSizes},
		{input: []int{1, 2, 3}},
	}

	for _, tc := range tt {
		err := s.Set(tc.input)
		if !errors.Is(err, tc.err) {
			t.Errorf("err should be %v; got: %v", tc.err, err)
		}
		if err == nil {
			sizes = tc.input
		}
		if got, want := s.Get(), sizes; !reflect.DeepEqual(got, want) {
			t.Errorf("%v: got: %v; want: %v", tc.input, got, want)
		}
	}
}
