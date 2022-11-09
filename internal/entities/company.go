package entities

import "sort"

type Company struct {
	ID        string     `json:"id"`
	Providers []Provider `json:"providers"`
}

func (c *Company) Sort() {
	sort.SliceStable(c.Providers, func(i, j int) bool {
		return c.Providers[i].Precedence < c.Providers[j].Precedence
	})
}
