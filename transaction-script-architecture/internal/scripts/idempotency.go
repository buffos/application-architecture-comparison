package scripts

import "transaction-script-architecture/internal/data"

func findIdempotentReturn(store *data.Store, command string, key string) (data.ReturnRequest, bool) {
	returnID, ok := store.Idempotency[command+":"+key]
	if !ok {
		return data.ReturnRequest{}, false
	}

	request, ok := store.Returns[returnID]
	return request, ok
}

func saveIdempotentReturn(store *data.Store, command string, key string, request data.ReturnRequest) {
	if store.Idempotency == nil {
		store.Idempotency = make(map[string]string)
	}
	store.Idempotency[command+":"+key] = request.ID
}
