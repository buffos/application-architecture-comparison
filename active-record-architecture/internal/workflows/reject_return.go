package workflows

import "active-record-architecture/internal/records"

// RejectReturn loads a return, records the review note, and persists the
// rejected state without applying inventory or refund side effects.
func RejectReturn(db *records.Database, returnID string, reviewNote string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	request, err := records.FindReturnRequest(db, returnID)
	if err != nil {
		return nil, err
	}
	if err := request.Reject(reviewNote); err != nil {
		return nil, err
	}
	return request, nil
}
