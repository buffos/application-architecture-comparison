package data

const (
	StockShortageRejectOrder    = "RejectOrder"
	StockShortageAllowBackorder = "AllowBackorder"
)

// StockRecord is a passive inventory record. Available stock is calculated by
// transaction scripts as OnHand minus Reserved.
type StockRecord struct {
	SKU              string
	OnHand           int
	Reserved         int
	ReorderThreshold int
}
