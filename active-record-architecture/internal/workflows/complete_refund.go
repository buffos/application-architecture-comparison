package workflows

import "active-record-architecture/internal/records"

// CompleteRefund loads an accepted return and invokes its reverse-side-effect
// operation.
func CompleteRefund(db *records.Database, returnID string, processedBy string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	request, err := records.FindReturnRequest(db, returnID)
	if err != nil {
		return nil, err
	}
	if err := request.CompleteRefund(processedBy); err != nil {
		return nil, err
	}
	return request, nil
}
