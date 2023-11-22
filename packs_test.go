package packs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/szkrstf/packs/mock"
)

func TestCalculate(t *testing.T) {
	tt := []struct {
		sizes []int
		items int
		packs map[int]int
	}{
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 0, packs: map[int]int{}},
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 1, packs: map[int]int{250: 1}},
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 250, packs: map[int]int{250: 1}},
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 251, packs: map[int]int{500: 1}},
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 501, packs: map[int]int{250: 1, 500: 1}},
		{sizes: []int{250, 500, 1000, 2000, 5000}, items: 12001, packs: map[int]int{250: 1, 2000: 1, 5000: 2}},
		{sizes: []int{23, 31, 53}, items: 152, packs: map[int]int{23: 2, 53: 2}},
		{sizes: []int{23, 31, 53}, items: 499_995, packs: map[int]int{23: 2, 53: 9433}},
		{sizes: []int{23, 31, 53}, items: 500_000, packs: map[int]int{23: 2, 31: 7, 53: 9429}},
		{sizes: []int{25, 50, 75}, items: 5_000_000_000, packs: map[int]int{50: 1, 75: 66666666}},
	}

	for _, tc := range tt {
		s := calculator{sizeStore: &mock.SizeStore{
			GetFn: func() []int { return tc.sizes },
		}}
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
