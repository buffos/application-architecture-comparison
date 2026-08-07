package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

// GetReturnRequest returns a defensive snapshot of one return request.
func GetReturnRequest(store *data.Store, returnID string) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}
	if returnID == "" {
		return data.ReturnRequest{}, ErrReturnIDRequired
	}

	request, ok := store.Returns[returnID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnNotFound
	}

	return cloneReturnRequest(request), nil
}

// ListReturnRequests lists defensive snapshots in deterministic ID order.
func ListReturnRequests(store *data.Store, status string) ([]data.ReturnRequest, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	requests := make([]data.ReturnRequest, 0, len(store.Returns))
	for _, request := range store.Returns {
		if status != "" && request.Status != status {
			continue
		}
		requests = append(requests, cloneReturnRequest(request))
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].ID < requests[j].ID
	})

	return requests, nil
}

func cloneReturnRequest(request data.ReturnRequest) data.ReturnRequest {
	request.Lines = append([]data.ReturnLine(nil), request.Lines...)
	return request
}
