package readmodel

import (
	"sort"

	"rules-engine-architecture/internal/engine"
)

type CustomerView struct {
	ID           string
	Name         string
	Tier         string
	InvoiceTerms bool
}

func ProjectCustomers(customers []engine.CustomerFact) []CustomerView {
	views := make([]CustomerView, 0, len(customers))
	for _, customer := range customers {
		views = append(views, CustomerView{
			ID:           customer.ID,
			Name:         customer.Name,
			Tier:         customer.Tier,
			InvoiceTerms: customer.InvoiceTerms,
		})
	}

	sort.Slice(views, func(left, right int) bool {
		return views[left].ID < views[right].ID
	})
	return views
}
