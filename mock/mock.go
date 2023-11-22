package mock

type Calculator struct {
	CalculateFn func(int) map[int]int
}

func (c *Calculator) Calculate(items int) map[int]int {
	return c.CalculateFn(items)
}

type SizeStore struct {
	GetFn func() []int
	SetFn func([]int) error
}

func (s *SizeStore) Get() []int {
	return s.GetFn()
}

func (s *SizeStore) Set(sizes []int) error {
	return s.SetFn(sizes)
}
