package workflows

import "active-record-architecture/internal/records"

// RequestReturn loads an order and asks it to create the passive return and
// refund records. The order itself is not changed by a return request.
func RequestReturn(db *records.Database, orderID string, lines []records.ReturnLine, reason string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	return order.RequestReturn(lines, reason)
}
