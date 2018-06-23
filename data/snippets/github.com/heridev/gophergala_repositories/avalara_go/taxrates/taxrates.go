package taxrates

type Rate struct {
	Rate float64 `json:"rate"`
	Name string  `json:"name"`
	Type string  `json:"type"`
}

type TaxResponse struct {
	TotalRate float64 `json:"totalRate"`
	Rates     []Rate  `json:"rates"`
}
