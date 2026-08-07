package workflows

import "active-record-architecture/internal/records"

// AcceptReturn loads a return and invokes its compressed acceptance
// operation. The return model persists the order, stock, refund, and return
// records involved in the reverse flow.
func AcceptReturn(db *records.Database, returnID string) (*records.ReturnRequest, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	request, err := records.FindReturnRequest(db, returnID)
	if err != nil {
		return nil, err
	}
	if err := request.Accept(); err != nil {
		return nil, err
	}
	return request, nil
}
