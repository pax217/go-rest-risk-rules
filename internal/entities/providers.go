package entities

type Provider struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Precedence int    `json:"precedence"`
	IsOn       bool   `json:"is_on"`
}
