package records

import "errors"

const (
	StockShortageRejectOrder    = "RejectOrder"
	StockShortageAllowBackorder = "AllowBackorder"
)

var (
	ErrProductSKURequired = errors.New("product sku is required")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductInactive    = errors.New("product is inactive")
)

// Product is an Active Record for the catalog table.
type Product struct {
	db *Database

	SKU                 string
	Name                string
	Category            string
	Active              bool
	UnitPrice           int
	ReturnWindowDays    int
	StockShortagePolicy string
}

// NewProduct creates a new, unsaved Product Active Record.
func NewProduct(db *Database, sku string, name string, category string, active bool, unitPrice int) *Product {
	return &Product{
		db:                  db,
		SKU:                 sku,
		Name:                name,
		Category:            category,
		Active:              active,
		UnitPrice:           unitPrice,
		StockShortagePolicy: StockShortageRejectOrder,
	}
}

// FindProduct loads a Product Active Record from the product table.
func FindProduct(db *Database, sku string) (*Product, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if sku == "" {
		return nil, ErrProductSKURequired
	}

	row, ok := db.products[sku]
	if !ok {
		return nil, ErrProductNotFound
	}

	return &Product{
		db:                  db,
		SKU:                 row.SKU,
		Name:                row.Name,
		Category:            row.Category,
		Active:              row.Active,
		UnitPrice:           row.UnitPrice,
		ReturnWindowDays:    row.ReturnWindowDays,
		StockShortagePolicy: row.StockShortagePolicy,
	}, nil
}

// Save writes the current Product Active Record to its table.
func (product *Product) Save() error {
	if product == nil || product.db == nil {
		return ErrDatabaseRequired
	}
	if product.SKU == "" {
		return ErrProductSKURequired
	}

	product.db.products[product.SKU] = productRow{
		SKU:                 product.SKU,
		Name:                product.Name,
		Category:            product.Category,
		Active:              product.Active,
		UnitPrice:           product.UnitPrice,
		ReturnWindowDays:    product.ReturnWindowDays,
		StockShortagePolicy: product.StockShortagePolicy,
	}
	return nil
}
