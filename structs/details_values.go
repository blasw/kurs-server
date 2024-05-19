package structs

type DetailValues struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
