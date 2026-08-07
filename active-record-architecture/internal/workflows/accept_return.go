package workflows

import (
	"time"

	"active-record-architecture/internal/records"
)

// AcceptReturn loads a return and invokes its compressed acceptance
// operation. The return model persists the order, stock, refund, and return
// records involved in the reverse flow.
func AcceptReturn(db *records.Database, returnID string, reviewedBy string, idempotencyKey string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	request, err := records.FindReturnRequest(db, returnID)
	if err != nil {
		return nil, err
	}
	return request.Accept(reviewedBy, idempotencyKey)
}

// AcceptReturnAt is the deterministic form used by tests and demonstrations.
func AcceptReturnAt(db *records.Database, returnID string, now time.Time, reviewedBy string, idempotencyKey string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	request, err := records.FindReturnRequest(db, returnID)
	if err != nil {
		return nil, err
	}
	return request.AcceptAt(now, reviewedBy, idempotencyKey)
}
