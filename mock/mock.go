package mock

type Calculator struct {
	CalculateFn func(int) map[int]int
}

func (c *Calculator) Calculate(items int) map[int]int {
	return c.CalculateFn(items)
}
