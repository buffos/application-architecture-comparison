package records

import "sort"

// GetReturnRequest is the named detail-query form of FindReturnRequest.
func GetReturnRequest(db *Database, id string) (*ReturnRequest, error) {
	return FindReturnRequest(db, id)
}

// ListReturnRequests returns reconstructed Active Record snapshots ordered by
// return ID. An empty status lists every request.
func ListReturnRequests(db *Database, status string) ([]*ReturnRequest, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	ids := make([]string, 0, len(db.returns))
	for id, row := range db.returns {
		if status != "" && row.Status != status {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	requests := make([]*ReturnRequest, 0, len(ids))
	for _, id := range ids {
		request, err := FindReturnRequest(db, id)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
